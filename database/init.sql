CREATE DATABASE IF NOT EXISTS users CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON users.* TO 'dtmnuser'@'%';

CREATE DATABASE IF NOT EXISTS api CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON api.* TO 'dtmnuser'@'%';

CREATE DATABASE IF NOT EXISTS integrations CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON integrations.* TO 'dtmnuser'@'%';

CREATE DATABASE IF NOT EXISTS pipelines CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON pipelines.* TO 'dtmnuser'@'%';

FLUSH PRIVILEGES;

SET GLOBAL max_allowed_packet=1073741824;
