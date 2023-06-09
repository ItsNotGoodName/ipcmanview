use anyhow::{Context, Result};

use chrono::{DateTime, Utc};
use sqlx::SqlitePool;

use crate::models::{Page, PageResult, ScanActive, ScanCompleted, ScanKind, ScanPending};
use crate::scan::{Scan, ScanActor, ScanCamera, ScanKindPending, ScanRange};

use super::NotFound;

impl Scan {
    pub(crate) async fn queue_all_db(pool: &SqlitePool, kind: ScanKindPending) -> Result<()> {
        let (range_start, range_end) = kind.range();
        let kind = ScanKind::from(kind);

        sqlx::query!(
            r#"
            REPLACE INTO pending_scans
            (
            camera_id,
            kind,
            range_start,
            range_end
            ) 
            SELECT id, ?, ?, ? from cameras
            "#,
            kind,
            range_start,
            range_end
        )
        .execute(pool)
        .await
        .context("Failed to insert into pending scans.")
        .map(|_| ())
    }

    pub(crate) async fn queue_db(
        pool: &SqlitePool,
        kind: ScanKindPending,
        camera_id: i64,
    ) -> Result<()> {
        let (range_start, range_end) = kind.range();
        let kind = ScanKind::from(kind);

        sqlx::query!(
            r#"
            REPLACE INTO pending_scans
            (
            camera_id,
            kind,
            range_start,
            range_end
            )
            VALUES (?, ?, ?, ?)
            "#,
            camera_id,
            kind,
            range_start,
            range_end
        )
        .execute(pool)
        .await
        .with_context(|| format!("Failed to insert into pending scans with camera {camera_id}."))
        .map(|_| ())
    }
}

impl ScanActor {
    pub(crate) async fn next(pool: &SqlitePool) -> Result<Option<Self>> {
        let mut pool = pool.begin().await?;

        // Create a actor from either pending_scans or pending_manual_scans, return if there is none
        let actor = if let Some(pending) = sqlx::query_as_unchecked!(ScanPending,
            "SELECT * FROM pending_scans WHERE camera_id NOT IN (SELECT camera_id FROM active_scans) LIMIT 1"
        ).fetch_optional(&mut pool).await? {
            // Delete pending scan
            sqlx::query!("DELETE FROM pending_scans WHERE id = ?", pending.id)
            .execute(&mut pool)
            .await?;

            // Create actor from pending scan kind
            match pending.kind {
                ScanKind::Full => ScanActor::full(pending.camera_id),
                ScanKind::Cursor => {
                    // Get scan camera
                    let scan_camera = sqlx::query_as_unchecked!(
                        ScanCamera,
                        "SELECT id, scan_cursor FROM cameras WHERE id = ?",
                        pending.camera_id
                    )
                    .fetch_one(&mut pool)
                    .await
                    .with_context(|| {
                        format!(
                            "Failed to find scan_cursor with camera {}.",
                            pending.camera_id
                        )
                    })?;

                    ScanActor::cursor(scan_camera)
                },
                ScanKind::Manual => {
                    ScanActor::manual(
                        pending.camera_id,
                        ScanRange {
                            start: pending.range_start,
                            end: pending.range_end,
                        },
                    )
                }
            }
        } else if let Some(completed) = sqlx::query_as_unchecked!(ScanCompleted,
            r#"
            SELECT * FROM completed_scans
            WHERE retry_pending = true
            AND camera_id NOT IN (SELECT camera_id FROM active_scans) LIMIT 1
            "#,
        ).fetch_optional(&mut pool).await? {
            sqlx::query!("UPDATE completed_scans SET retry_pending = false, can_retry = false WHERE id = ?", completed.id)
            .execute(&mut pool)
            .await?;

            ScanActor::from(completed)
        } else {
            return Ok(None);
        };

        // Insert actor into active scans
        sqlx::query!(
            r#"
            INSERT INTO active_scans
            (
            camera_id,
            kind,
            range_start,
            range_end,
            started_at,
            range_cursor
            )
            VALUES
            (?, ?, ?, ?, ?, ?)
            "#,
            actor.camera_id,
            actor.kind,
            actor.range.start,
            actor.range.end,
            actor.started_at,
            actor.range.end,
        )
        .execute(&mut pool)
        .await
        .with_context(|| {
            format!(
                "Failed to create active scan with camera {}.",
                actor.camera_id
            )
        })?;

        pool.commit().await?;

        Ok(Some(actor))
    }

    pub(crate) async fn update_status(
        &self,
        pool: &SqlitePool,
        range_cursor: DateTime<Utc>,
        percent: f64,
        upserted: i64,
        deleted: i64,
    ) -> Result<()> {
        sqlx::query!(
            r#"
            UPDATE active_scans SET
            range_cursor = ?,
            percent = ?,
            upserted = ?,
            deleted = ?
            WHERE camera_id = ?
            "#,
            range_cursor,
            percent,
            upserted,
            deleted,
            self.camera_id,
        )
        .execute(pool)
        .await
        .with_context(|| {
            format!(
                "Failed to update status on active scan with camera {}.",
                self.camera_id
            )
        })
        .map(|_| ())
    }

    pub(crate) async fn end(self, pool: &SqlitePool) -> Result<()> {
        let mut pool = pool.begin().await?;

        // Save scan actor to completed_scans
        if self.should_save() {
            let duration = self.duration();
            let success = self.success();
            let can_retry = self.can_retry();
            sqlx::query!(
                r#"
                INSERT INTO completed_scans 
                (
                camera_id,
                kind,
                range_start,
                range_end,
                started_at,
                range_cursor,
                deleted,
                upserted,
                percent,
                duration,
                success,
                can_retry,
                error
                )
                SELECT
                camera_id,
                kind,
                range_start,
                range_end,
                started_at,
                range_cursor,
                deleted,
                upserted,
                percent,
                ?,
                ?,
                ?,
                ?
                FROM active_scans WHERE camera_id = ?
                "#,
                duration,
                success,
                can_retry,
                self.error,
                self.camera_id
            )
            .execute(&mut pool)
            .await
            .with_context(|| {
                format!(
                    "Failed to insert into completed scans with camera {}.",
                    self.camera_id
                )
            })?;
        }

        // Delete scan actor from active_scans
        sqlx::query!(
            "DELETE FROM active_scans WHERE camera_id = ?",
            self.camera_id
        )
        .execute(&mut pool)
        .await
        .with_context(|| {
            format!(
                "Failed to delete active scan with camera {}.",
                self.camera_id
            )
        })?;

        // Update camera scan cursor
        if let Some(scan_cursor) = self.should_update_scan_cursor() {
            sqlx::query!(
                "UPDATE cameras SET scan_cursor = ?1 WHERE id = ?2 AND scan_cursor < ?1",
                scan_cursor,
                self.camera_id,
            )
            .execute(&mut pool)
            .await
            .with_context(|| {
                format!(
                    "Failed to update scan cursor with camera {}.",
                    self.camera_id,
                )
            })?;
        }

        pool.commit().await?;

        Ok(())
    }
}

impl ScanActive {
    pub(crate) async fn clear(pool: &SqlitePool) -> Result<()> {
        sqlx::query!("DELETE FROM active_scans")
            .execute(pool)
            .await
            .context("Failed to delete active scans.")
            .map(|_| ())
    }

    pub async fn list(pool: &SqlitePool) -> Result<Vec<Self>> {
        sqlx::query_as_unchecked!(
            Self,
            r#"
            SELECT * FROM active_scans
            "#
        )
        .fetch_all(pool)
        .await
        .context("Failed to list active scans.")
    }
}

impl ScanCompleted {
    // TODO: add filters
    pub async fn list(pool: &SqlitePool, page: Page) -> Result<PageResult<Self>> {
        let offset = page.offset();
        let limit = page.limit();
        let items = sqlx::query_as_unchecked!(
            Self,
            r#"
            SELECT *
            FROM completed_scans
            ORDER BY started_at DESC
            LIMIT ?
            OFFSET ?
            "#,
            limit,
            offset,
        )
        .fetch_all(pool)
        .await
        .context("Failed to list completed scans.")?;
        let total_items = sqlx::query!(r#"SELECT COUNT(id) AS count FROM completed_scans"#)
            .fetch_one(pool)
            .await
            .with_context(|| format!("Failed to count completed scans."))
            .map(|c| c.count)?;

        Ok(page.new_result(items, total_items))
    }

    pub async fn find(pool: &SqlitePool, id: i64) -> Result<Self> {
        sqlx::query_as_unchecked!(
            Self,
            r#"
            SELECT *
            FROM completed_scans
            WHERE id = ?
            "#,
            id
        )
        .fetch_optional(pool)
        .await
        .with_context(|| format!("Failed to find completed scan {id}."))?
        .ok_or(NotFound)
        .with_context(|| format!("Failed to find completed scan {id}."))
    }

    pub(crate) async fn retry_db(pool: &SqlitePool, id: i64) -> Result<()> {
        sqlx::query!(
            "UPDATE completed_scans SET retry_pending = true WHERE id = ? AND can_retry = true",
            id
        )
        .execute(pool)
        .await
        .with_context(|| format!("Failed to retry completed scan {id}."))
        .map(NotFound::check_query)?
        .with_context(|| format!("Failed to retry completed scan {id}."))
    }
}

impl ScanPending {
    pub async fn list(pool: &SqlitePool) -> Result<Vec<Self>> {
        sqlx::query_as_unchecked!(Self, "SELECT * FROM pending_scans")
            .fetch_all(pool)
            .await
            .context("Failed to list pending scans.")
    }
}
