// Define an environment named "local"
env "local" {
  // Declare where the schema definition resides.
  // Also supported: ["file://multi.hcl", "file://schema.hcl"].
  src = "file://internal/migrations/schema.sql"

  // Define the URL of the database which is managed
  // in this environment.
  url = "sqlite://sqlite.db?_fk=1"

  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  dev = "sqlite://file?mode=memory&_fk=1"

  migration {
    // URL where the migration directory resides.
    dir = "file://internal/migrations/sql"
    // An optional format of the migration directory:
    // atlas (default) | flyway | liquibase | goose | golang-migrate | dbmate
    format = goose
  }
}
