import axios from "axios";

const INTEGRATIONS_API_URL = "/integration-api/";

export const INTEGRATION_API_METHOD_POST = "post";
export const INTEGRATION_API_METHOD_GET = "get";
export const INTEGRATION_API_METHOD_PUT = "put";
export const INTEGRATION_API_METHOD_PATCH = "patch";
export const INTEGRATION_API_METHOD_DELETE = "delete";
export const INTEGRATION_API_METHOD_COPY = "copy";
export const INTEGRATION_API_METHOD_HEAD = "head";
export const INTEGRATION_API_METHOD_OPTIONS = "options";
export const INTEGRATION_API_METHOD_LINK = "link";
export const INTEGRATION_API_METHOD_UNLINK = "unlink";
export const INTEGRATION_API_METHOD_PURGE = "purge";
export const INTEGRATION_API_METHOD_LOCK = "lock";
export const INTEGRATION_API_METHOD_UNLOCK = "unlock";
export const INTEGRATION_API_METHOD_VIEW = "view";
export const INTEGRATION_API_METHOD_PROPFIND = "propfind";

export const INTEGRATION_API_METHODS = [
    INTEGRATION_API_METHOD_POST,
    INTEGRATION_API_METHOD_GET,
    INTEGRATION_API_METHOD_PUT,
    INTEGRATION_API_METHOD_PATCH,
    INTEGRATION_API_METHOD_DELETE,
    INTEGRATION_API_METHOD_COPY,
    INTEGRATION_API_METHOD_HEAD,
    INTEGRATION_API_METHOD_OPTIONS,
    INTEGRATION_API_METHOD_LINK,
    INTEGRATION_API_METHOD_UNLINK,
    INTEGRATION_API_METHOD_PURGE,
    INTEGRATION_API_METHOD_LOCK,
    INTEGRATION_API_METHOD_UNLOCK,
    INTEGRATION_API_METHOD_VIEW,
    INTEGRATION_API_METHOD_PROPFIND,
];

export const INTEGRATION_IO_TYPE_READ_WRITE = "read-write";
export const INTEGRATION_IO_TYPE_READ = "read";
export const INTEGRATION_IO_TYPE_WRITE = "write";
export const INTEGRATION_IO_TYPE_ALL = "all";

export const INTEGRATION_IO_TYPES = [
    INTEGRATION_IO_TYPE_READ_WRITE,
    INTEGRATION_IO_TYPE_WRITE,
    INTEGRATION_IO_TYPE_READ,
];

export const INTEGRATION_IO_TYPES_FORM = [
    INTEGRATION_IO_TYPE_ALL,
    INTEGRATION_IO_TYPE_READ_WRITE,
    INTEGRATION_IO_TYPE_WRITE,
    INTEGRATION_IO_TYPE_READ,
];

export const INTEGRATION_TYPE_API = "api";
export const INTEGRATION_TYPE_SMS = "sms";
export const INTEGRATION_TYPE_WHATSAPP = "whatsapp"
export const INTEGRATION_TYPE_EMAIL = "email";
export const INTEGRATION_TYPE_SLACK = "slack";
export const INTEGRATION_TYPE_JIRA = "jira";
export const INTEGRATION_TYPE_INCIDENT_IO = "incidentio";
export const INTEGRATION_TYPE_TABLEAU = "tableau";
export const INTEGRATION_TYPE_HUBSPOT = "hubspot";
export const INTEGRATION_TYPE_SALESFORCE = "salesforce";
export const INTEGRATION_TYPE_GOOGLE_SHEETS = "google-sheets";
export const INTEGRATION_TYPE_OPSGENIE = "opsgenie";
export const INTEGRATION_TYPE_JENKINS = "jenkins";

export const INTEGRATION_TYPE_SQL = "sql";
export const INTEGRATION_TYPE_MYSQL = "mysql";
export const INTEGRATION_TYPE_SNOWFLAKE = "snowflake";
export const INTEGRATION_TYPE_POSTGRESQL = "postgresql";
export const INTEGRATION_TYPE_AWS_RDS = "aws-rds";
export const INTEGRATION_TYPE_GOOGLE_CLOUD_SQL = "google-cloud-sql";
export const INTEGRATION_TYPE_GOOGLE_BIG_QUERY = "google-bigquery";
export const INTEGRATION_TYPE_PLANET_SCALE = "planet-scale";
export const INTEGRATION_TYPE_IMMUTA = "immuta";
export const INTEGRATION_TYPE_MICROSOFT_AZURE_SQL = "microsoft-azure-sql";
export const INTEGRATION_TYPE_ELASTICSEARCH = "elasticsearch";
export const INTEGRATION_TYPE_REDSHIFT = "redshift";
export const INTEGRATION_TYPE_CLICKHOUSE = "clickhouse";

export const INTEGRATION_TYPES = [
    INTEGRATION_TYPE_API,
    INTEGRATION_TYPE_SLACK,
    INTEGRATION_TYPE_JIRA,
    INTEGRATION_TYPE_EMAIL,
    INTEGRATION_TYPE_INCIDENT_IO,
    INTEGRATION_TYPE_SMS,
    INTEGRATION_TYPE_WHATSAPP,
    INTEGRATION_TYPE_TABLEAU,
    INTEGRATION_TYPE_HUBSPOT,
    INTEGRATION_TYPE_SALESFORCE,
    INTEGRATION_TYPE_GOOGLE_SHEETS,
    INTEGRATION_TYPE_OPSGENIE,
    INTEGRATION_TYPE_JENKINS,
    INTEGRATION_TYPE_SQL,
];

export const INTEGRATION_TYPES_PER_IO = {
    [INTEGRATION_IO_TYPE_READ_WRITE]: [
        INTEGRATION_TYPE_API,
        INTEGRATION_TYPE_SQL,
    ],
    [INTEGRATION_IO_TYPE_READ]: [
        INTEGRATION_TYPE_SQL,
    ],
    [INTEGRATION_IO_TYPE_WRITE]: [
        INTEGRATION_TYPE_SLACK,
        INTEGRATION_TYPE_JIRA,
        INTEGRATION_TYPE_EMAIL,
        INTEGRATION_TYPE_INCIDENT_IO,
        INTEGRATION_TYPE_SMS,
        INTEGRATION_TYPE_WHATSAPP,
        INTEGRATION_TYPE_TABLEAU,
        INTEGRATION_TYPE_HUBSPOT,
        INTEGRATION_TYPE_SALESFORCE,
        INTEGRATION_TYPE_GOOGLE_SHEETS,
        INTEGRATION_TYPE_OPSGENIE,
        INTEGRATION_TYPE_JENKINS,
    ],
};

export const INTEGRATION_TYPES_PER_IO_FORM = {
    [INTEGRATION_IO_TYPE_READ_WRITE]: [
        INTEGRATION_TYPE_API,
        INTEGRATION_TYPE_MYSQL,
        INTEGRATION_TYPE_POSTGRESQL,
        INTEGRATION_TYPE_SNOWFLAKE,
        INTEGRATION_TYPE_REDSHIFT,
        INTEGRATION_TYPE_AWS_RDS,
        INTEGRATION_TYPE_GOOGLE_CLOUD_SQL,
        INTEGRATION_TYPE_GOOGLE_BIG_QUERY,
        INTEGRATION_TYPE_MICROSOFT_AZURE_SQL,
        INTEGRATION_TYPE_PLANET_SCALE,
        INTEGRATION_TYPE_IMMUTA,
        INTEGRATION_TYPE_CLICKHOUSE,
    ],
    [INTEGRATION_IO_TYPE_WRITE]: [
        INTEGRATION_TYPE_SLACK,
        INTEGRATION_TYPE_JIRA,
        INTEGRATION_TYPE_EMAIL,
        INTEGRATION_TYPE_INCIDENT_IO,
        INTEGRATION_TYPE_SMS,
        INTEGRATION_TYPE_WHATSAPP,
        INTEGRATION_TYPE_TABLEAU,
        INTEGRATION_TYPE_HUBSPOT,
        INTEGRATION_TYPE_SALESFORCE,
        INTEGRATION_TYPE_GOOGLE_SHEETS,
        INTEGRATION_TYPE_OPSGENIE,
        INTEGRATION_TYPE_JENKINS,
    ],
    [INTEGRATION_IO_TYPE_READ]: [
        INTEGRATION_TYPE_ELASTICSEARCH,
    ],
};

export const INTEGRATION_SQL_TYPES_PER_IO = {
    [INTEGRATION_IO_TYPE_READ_WRITE]: [
        INTEGRATION_TYPE_MYSQL,
        INTEGRATION_TYPE_POSTGRESQL,
        INTEGRATION_TYPE_SNOWFLAKE,
        INTEGRATION_TYPE_REDSHIFT,
        INTEGRATION_TYPE_AWS_RDS,
        INTEGRATION_TYPE_GOOGLE_CLOUD_SQL,
        INTEGRATION_TYPE_GOOGLE_BIG_QUERY,
        INTEGRATION_TYPE_MICROSOFT_AZURE_SQL,
        INTEGRATION_TYPE_PLANET_SCALE,
        INTEGRATION_TYPE_IMMUTA,
        INTEGRATION_TYPE_CLICKHOUSE,
    ],
    [INTEGRATION_IO_TYPE_READ]: [
        INTEGRATION_TYPE_ELASTICSEARCH,
    ],
    [INTEGRATION_IO_TYPE_WRITE]: [],
};

export const INTEGRATION_TYPE_TO_HUMAN = {
    [INTEGRATION_TYPE_API]: "API",
    [INTEGRATION_TYPE_SLACK]: "Slack",
    [INTEGRATION_TYPE_JIRA]: "Atlassian Jira",
    [INTEGRATION_TYPE_EMAIL]: "E-mail",
    [INTEGRATION_TYPE_INCIDENT_IO]: "Incident.io",
    [INTEGRATION_TYPE_SMS]: "SMS (by Twilio)",
    [INTEGRATION_TYPE_WHATSAPP]: "WhatsApp (by Twilio)",
    [INTEGRATION_TYPE_TABLEAU]: "Tableau",
    [INTEGRATION_TYPE_HUBSPOT]: "Hubspot",
    [INTEGRATION_TYPE_SALESFORCE]: "Salesforce",
    [INTEGRATION_TYPE_GOOGLE_SHEETS]: "Google Sheets",
    [INTEGRATION_TYPE_OPSGENIE]: "Opsgenie",
    [INTEGRATION_TYPE_JENKINS]: "Jenkins",
    [INTEGRATION_TYPE_SQL]: "SQL",
    [INTEGRATION_TYPE_SNOWFLAKE]: "Snowflake",
    [INTEGRATION_TYPE_REDSHIFT]: "Amazon Redshift",
    [INTEGRATION_TYPE_MYSQL]: "MySQL",
    [INTEGRATION_TYPE_POSTGRESQL]: "PostgreSQL",
    [INTEGRATION_TYPE_AWS_RDS]: "AWS RDS",
    [INTEGRATION_TYPE_GOOGLE_CLOUD_SQL]: "Google Cloud SQL",
    [INTEGRATION_TYPE_GOOGLE_BIG_QUERY]: "Google Big Query",
    [INTEGRATION_TYPE_MICROSOFT_AZURE_SQL]: "Microsoft Azure SQL",
    [INTEGRATION_TYPE_PLANET_SCALE]: "PlanetScale",
    [INTEGRATION_TYPE_IMMUTA]: "Immuta",
    [INTEGRATION_TYPE_ELASTICSEARCH]: "ElasticSearch",
    [INTEGRATION_TYPE_CLICKHOUSE]: "ClickHouse",
};

export const INTEGRATION_TYPE_API_AUTH_TYPE_NONE = "None";
export const INTEGRATION_TYPE_API_AUTH_TYPE_BASIC = "Basic";
export const INTEGRATION_TYPE_API_AUTH_TYPE_BEARER = "Bearer";
export const INTEGRATION_TYPE_API_AUTH_TYPE_HEADER = "Header";

export const INTEGRATION_TYPE_API_AUTH_TYPES = [
    INTEGRATION_TYPE_API_AUTH_TYPE_NONE,
    INTEGRATION_TYPE_API_AUTH_TYPE_BASIC,
    INTEGRATION_TYPE_API_AUTH_TYPE_BEARER,
    INTEGRATION_TYPE_API_AUTH_TYPE_HEADER,
];

export const SQL_INTEGRATION_TYPES = [
    INTEGRATION_TYPE_SNOWFLAKE,
    INTEGRATION_TYPE_REDSHIFT,
    INTEGRATION_TYPE_MYSQL,
    INTEGRATION_TYPE_POSTGRESQL,
    INTEGRATION_TYPE_AWS_RDS,
    INTEGRATION_TYPE_GOOGLE_CLOUD_SQL,
    INTEGRATION_TYPE_GOOGLE_BIG_QUERY,
    INTEGRATION_TYPE_MICROSOFT_AZURE_SQL,
    INTEGRATION_TYPE_PLANET_SCALE,
    INTEGRATION_TYPE_IMMUTA,
    INTEGRATION_TYPE_ELASTICSEARCH,
    INTEGRATION_TYPE_CLICKHOUSE,
];

export const SQL_INTEGRATION_TYPE_PORT_MAP = {
    [INTEGRATION_TYPE_MYSQL]: 3306,
    [INTEGRATION_TYPE_SNOWFLAKE]: 443,
    [INTEGRATION_TYPE_POSTGRESQL]: 5432,
    [INTEGRATION_TYPE_AWS_RDS]: 3306,
    [INTEGRATION_TYPE_GOOGLE_CLOUD_SQL]: 3306,
    [INTEGRATION_TYPE_GOOGLE_BIG_QUERY]: null,
    [INTEGRATION_TYPE_MICROSOFT_AZURE_SQL]: 3306,
    [INTEGRATION_TYPE_PLANET_SCALE]: 3306,
    [INTEGRATION_TYPE_IMMUTA]: 5432,
    [INTEGRATION_TYPE_ELASTICSEARCH]: 9200,
    [INTEGRATION_TYPE_REDSHIFT]: 5439,
};

export const SQL_INTEGRATION_CONNECTION_TYPE_DIRECT = "direct";
export const SQL_INTEGRATION_CONNECTION_TYPE_SSH = "ssh";

export const SQL_INTEGRATION_CONNECTION_TYPES = [
    SQL_INTEGRATION_CONNECTION_TYPE_DIRECT,
    SQL_INTEGRATION_CONNECTION_TYPE_SSH,
];

export const SQL_INTEGRATION_SSL_MODES = {
    "Off": false,
    "On": true,
};

export const ELASTICSEARCH_VERSIONS = [6, 7, 8];

class IntegrationService {
    getIntegrationsByOrganization = async(uuid, ioType = "all") => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'organization/' + uuid + '/integrations/' + ioType,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getIntegration(uuid, type, sqlType = false) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                INTEGRATIONS_API_URL + type + '/' + (sqlType !== false ? sqlType + '/' : '' ) + uuid,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    deleteIntegration(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                INTEGRATIONS_API_URL + 'integration/' + uuid + '/delete',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updateIntegration(uuid, type, data, sqlType = false) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                INTEGRATIONS_API_URL + type + '/' + (sqlType !== false ? sqlType + '/' : '' ) + uuid,
                data,
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then((response) => {
                return response.data;
            });
    }

    addIntegration(type, data, sqlType = false) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                INTEGRATIONS_API_URL + type + (sqlType !== false ? '/' + sqlType : ''),
                data,
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then((response) => {
                return response.data;
            });
    }

    confirmIntegration(type, code, uuid) {
        var token = localStorage.getItem("token")

        return axios
            .post(
                INTEGRATIONS_API_URL + type + '/' + uuid + '/confirm',
                { code },
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getSlackAuthorizations = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'organization/' + uuid + '/slack/authorizations',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    createSlackAuthorization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'organization/' + uuid + '/slack/authorization',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    reauthorizeSlackAuthorization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'slack/authorization/' + uuid + '/reauthorize',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updateSlackAuthorization = async(name, uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'slack/authorization/' + uuid,
                {
                    name
                },
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getJiraAuthorizations = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'organization/' + uuid + '/jira/authorizations',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getJiraAuthorization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'jira/authorization/' + uuid,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    createJiraAuthorization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'organization/' + uuid + '/jira/authorization',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updateJiraAuthorization = async(name, resource_id, uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'jira/authorization/' + uuid,
                {
                    name,
                    resource_id
                },
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getIncidentIoIntegrationSeverities = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'incidentio/' + uuid + '/severities',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getHubspotAuthorizations(uuid) {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'organization/' + uuid + '/hubspot/authorizations',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getHubspotAuthorization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'hubspot/authorization/' + uuid,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    createHubspotAuthorization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'organization/' + uuid + '/hubspot/authorization',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updateHubspotAuthorization = async(name, uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'hubspot/authorization/' + uuid,
                {
                    name,
                },
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getSalesforceAuthorizations(uuid) {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'organization/' + uuid + '/salesforce/authorizations',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getSalesforceAuthorization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'salesforce/authorization/' + uuid,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    createSalesforceAuthorization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'organization/' + uuid + '/salesforce/authorization',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updateSalesforceAuthorization = async(name, uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'salesforce/authorization/' + uuid,
                {
                    name,
                },
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    testExistingIntegration(integrationUuid, type, data) {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'integration/sql/' + integrationUuid + '/test',
                data,
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    testNewIntegration(type, data) {
        var token = localStorage.getItem("token");

        return axios
            .post(
                INTEGRATIONS_API_URL + 'integration/sql/test/' + type,
                data,
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getDatabases(integrationUuid) {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'integration/sql/' + integrationUuid + '/dbs',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getTables(integrationUuid, db) {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'integration/sql/' + integrationUuid + '/db/' + db + '/tables',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    describeTable(integrationUuid, db, table) {
        var token = localStorage.getItem("token");

        return axios
            .get(
                INTEGRATIONS_API_URL + 'integration/sql/' + integrationUuid + '/db/' + db + '/table/' + table,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }
}

export default new IntegrationService();
