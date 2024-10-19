import axios from "axios";
import {Buffer} from 'buffer';

const PIPELINE_API_URL = "/pipeline-api/";

export const PIPELINE_ROOT_FOLDER = "root";

export const PIPELINE_TYPE_GENERIC = "generic";
export const PIPELINE_TYPE_METRIC = "metric";

export const PIPELINE_PAGE_PREVIEW = "preview";
export const PIPELINE_PAGE_DETAILS = "details";
export const PIPELINE_PAGE_STATS = "stats";
export const PIPELINE_PAGE_LOGS = "logs";
export const PIPELINE_PAGE_TRIGGERS = "triggers";

export const PIPELINE_PAGES = {
    PIPELINE_PAGE_PREVIEW,
    PIPELINE_PAGE_DETAILS,
    PIPELINE_PAGE_STATS,
    PIPELINE_PAGE_LOGS,
    PIPELINE_PAGE_TRIGGERS,
};

class PipelineService {
    getPipelinesByOrganization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                PIPELINE_API_URL + 'organization/' + uuid + '/pipelines',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getDashboardByOrganization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                PIPELINE_API_URL + 'organization/' + uuid + '/dashboard',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getNewGroupedItemsByOrganization = async(uuid, type, groupBy) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                PIPELINE_API_URL + 'organization/' + uuid + '/new-grouped-items/' + type + '/' + groupBy,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getRunsPerMonthByOrganization = async(uuid, type) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                PIPELINE_API_URL + 'organization/' + uuid + '/runs_per_month/' + type,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    searchPipelines = async(uuid, string) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                PIPELINE_API_URL + 'organization/' + uuid + '/pipelines/search/' + string,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getPipelinesByOrganizationAndFolder = async(organizationUuid, folderUuid = null) => {
        var token = localStorage.getItem("token");

        var url = folderUuid === null
            ? PIPELINE_API_URL + 'organization/' + organizationUuid + '/root_folder/pipelines'
            : PIPELINE_API_URL + 'organization/' + organizationUuid + '/folder/' + folderUuid + '/pipelines'
        ;

        return axios
            .get(
                url,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    deletePipeline(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + uuid + '/delete',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    copyPipeline(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + uuid + '/clone',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then((response) => {
                return response.data;
            })
            .catch((error) => {
                return error;
            });
    }

    togglePipeline(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + uuid + '/toggle',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    runPipeline(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + uuid + '/run',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getPipeline(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                PIPELINE_API_URL + 'pipeline/' + uuid,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getPipelineRunResults(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                PIPELINE_API_URL + 'pipeline/' + uuid + '/run',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updatePipeline(uuid, name, folder_uuid, elements_layout = [], schedule = null) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + uuid,
                { name, folder_uuid, schedule, elements_layout: JSON.stringify(elements_layout) },
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    addPipeline(name, folder_uuid, organization_uuid, elements_layout = [], type = PIPELINE_TYPE_GENERIC, schedule = null) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline',
                { name, schedule, folder_uuid, type, organization_uuid, elements_layout: JSON.stringify(elements_layout)},
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then((response) => {
                return response.data;
            });
    }

    getPipelinePreview(uuid, asTemplate = false) {
        var token = localStorage.getItem("token");
        var url = PIPELINE_API_URL + 'pipeline/' + uuid + '/preview';

        if (asTemplate === true) {
            url = url + "?asTemplate=1";
        }

        return axios
            .get(
                url,
                { responseType: 'arraybuffer' },
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then(response => Buffer.from(response.data, 'binary').toString('base64'));
    }

    updatePipelinePreview(uuid, preview) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + uuid + '/preview',
                preview,
                { headers: { Authorization: 'Bearer ' + token, "content-type": preview.type } }
            );
    }
}

export default new PipelineService();
