/* pgmigrate-encoding: utf-8 */

CREATE TABLE pdf_jobs (
    id SERIAL PRIMARY KEY,
    file_path VARCHAR(128)
)