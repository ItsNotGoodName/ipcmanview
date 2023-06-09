use anyhow::Result;
use sqlx::SqlitePool;

use crate::ipc::{IpcDetail, IpcLicenses, IpcManager, IpcSoftware, IpcStore};
use crate::models::{
    Camera, CameraDetail, CameraFile, CameraFileQuery, CameraFileQueryCursor,
    CameraFileQueryFilter, CameraFileQueryResult, CameraLicense, CameraScanResult, CameraShow,
    CameraSoftware, CreateCameraRequest, ScanCompleted, UpdateCameraRequest,
};
use crate::scan::{Scan, ScanActor, ScanKindPending};

// -------------------- Camera

impl CreateCameraRequest {
    pub async fn create(self, pool: &SqlitePool, store: &IpcStore) -> Result<i64> {
        // Create in database
        let id = self.create_db(pool).await?;
        // Refresh in store
        store.refresh(id).await?;
        // Get from store and refresh in database
        store.get(id).await?.refresh(pool).await.ok();
        // Queue a full scan
        Scan::queue(pool, store, id, ScanKindPending::Full).await?;

        Ok(id)
    }
}

impl UpdateCameraRequest {
    pub async fn update(self, pool: &SqlitePool, store: &IpcStore) -> Result<()> {
        let id = self.id;
        // Update in database
        self.update_db(pool).await?;
        // Refresh in store
        store.refresh(id).await?;
        // Get from store and refresh in database
        store.get(id).await?.refresh(pool).await.ok();

        Ok(())
    }
}

impl Camera {
    pub async fn delete(pool: &SqlitePool, store: &IpcStore, id: i64) -> Result<()> {
        // Delete in database
        Self::delete_db(pool, id).await?;
        // Refresh in store
        store.refresh(id).await?;

        Ok(())
    }
}

impl CameraShow {
    // TODO: make this into a single query inside of db crate
    pub async fn find(pool: &SqlitePool, id: i64) -> Result<Self> {
        let detail = CameraDetail::find(pool, id).await?;
        let software = CameraSoftware::find(pool, id).await?;
        let licenses = CameraLicense::list(pool, id).await?;
        let camera = Camera::find(pool, id).await?;

        let file_total =
            CameraFile::total(pool, &CameraFileQueryFilter::new().camera_ids(vec![id])).await?;

        Ok(CameraShow {
            id: camera.id,
            ip: camera.ip,
            username: camera.username,
            refreshed_at: camera.refreshed_at,
            created_at: camera.created_at,
            detail,
            software,
            licenses,
            file_total,
        })
    }
}

impl IpcManager {
    pub async fn refresh(&self, pool: &SqlitePool) -> Result<()> {
        self.refresh_detail(pool).await?;
        self.refresh_licenses(pool).await?;
        self.refresh_software(pool).await?;
        Camera::update_refreshed_at(pool, self.id).await
    }

    pub async fn refresh_detail(&self, pool: &SqlitePool) -> Result<()> {
        IpcDetail::get(&self).await?.save(pool, self.id).await
    }

    pub async fn refresh_licenses(&self, pool: &SqlitePool) -> Result<()> {
        IpcLicenses::get(&self).await?.save(pool, self.id).await
    }

    pub async fn refresh_software(&self, pool: &SqlitePool) -> Result<()> {
        IpcSoftware::get(&self).await?.save(pool, self.id).await
    }
}

impl CameraFile {
    pub async fn query(
        pool: &SqlitePool,
        store: &IpcStore,
        query: CameraFileQuery<'_>,
    ) -> Result<CameraFileQueryResult> {
        // Cursor scan when no cursor is supplied
        if let CameraFileQueryCursor::None = query.cursor {
            Scan::queue_all(pool, store, ScanKindPending::Cursor).await?;
        }

        Self::query_db(pool, query).await
    }
}

// -------------------- Scan

impl Scan {
    pub async fn queue(
        pool: &SqlitePool,
        store: &IpcStore,
        camera_id: i64,
        kind: ScanKindPending,
    ) -> Result<()> {
        Self::queue_db(pool, kind, camera_id).await?;
        Self::run_pending(pool, store).await;
        Ok(())
    }

    pub async fn queue_all(
        pool: &SqlitePool,
        store: &IpcStore,
        kind: ScanKindPending,
    ) -> Result<()> {
        Self::queue_all_db(pool, kind).await?;
        Self::run_pending(pool, store).await;
        Ok(())
    }

    // TODO: return database access errors
    pub async fn run_pending(pool: &SqlitePool, store: &IpcStore) {
        // Get a pending scan
        let first_handle = if let Ok(Some(s)) = ScanActor::next(pool).await {
            s
        } else {
            return;
        };

        // Get rest of the pending scans
        let mut handles = vec![first_handle];
        loop {
            match ScanActor::next(&pool).await {
                Ok(Some(handle)) => handles.push(handle),
                Ok(None) | Err(_) => break,
            }
        }

        // Start worker for each scan
        for handle in handles {
            let pool = pool.clone();
            let store = store.clone();
            tokio::spawn(async move {
                // Run pending scan
                if let Err(err) = handle.run(&pool, &store).await {
                    tracing::error!("{err:?}");
                }
                // Check for more scans and run them or exit
                loop {
                    match ScanActor::next(&pool).await {
                        Ok(Some(handle)) => {
                            if let Err(err) = handle.run(&pool, &store).await {
                                tracing::error!("{err:?}");
                            }
                        }
                        Ok(None) => return,
                        Err(err) => {
                            tracing::error!("{err:?}");
                            return;
                        }
                    }
                }
            });
        }
    }
}

impl ScanCompleted {
    pub async fn retry(pool: &SqlitePool, store: &IpcStore, id: i64) -> Result<()> {
        Self::retry_db(pool, id).await?;
        Scan::run_pending(pool, store).await;
        Ok(())
    }
}

impl ScanActor {
    async fn runner(&self, pool: &SqlitePool, man: &IpcManager) -> Result<()> {
        let mut res = CameraScanResult::default();
        for (range, percent) in self.range.iter() {
            res += man.scan_files(pool, range.start, range.end).await?;
            self.update_status(
                pool,
                range.start,
                percent,
                res.upserted as i64,
                res.deleted as i64,
            )
            .await?
        }

        Ok(())
    }

    async fn run(mut self, pool: &SqlitePool, store: &IpcStore) -> Result<()> {
        // Get manager
        let man = store.get(self.camera_id).await?;

        // Run scan
        let res = self.runner(pool, &man).await;
        if let Err(ref err) = res {
            self.error = format!("{:?}", err)
        }

        // End scan
        self.end(pool).await?;

        res
    }
}
