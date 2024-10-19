import axios from "axios";

const PIPELINE_API_URL = "/pipeline-api/";

class PipelineService {
    getFoldersByOrganizationAndFolder = async(organizationUuid, folderUuid = null) => {
        var token = localStorage.getItem("token");

        var url = folderUuid === null
            ? PIPELINE_API_URL + 'organization/' + organizationUuid + '/folders'
            : PIPELINE_API_URL + 'organization/' + organizationUuid + '/folder/' + folderUuid + '/folders'
        ;

        return axios
            .get(
                url,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    deleteFolder(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'folder/' + uuid + '/delete',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getFolder(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                PIPELINE_API_URL + 'folder/' + uuid,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updateFolder(uuid, name, parent_uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'folder/' + uuid,
                { name, parent_uuid },
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    addFolder(name, parent_uuid, organization_uuid, type) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'folder',
                { name, organization_uuid, parent_uuid, type},
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then((response) => {
                return response.data;
            });
    }
}

export default new PipelineService();
