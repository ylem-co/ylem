import axios from "axios";

const PIPELINE_API_URL = "/pipeline-api/";

class EnvVariableService {
    getVariablesByOrganization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                PIPELINE_API_URL + 'organization/' + uuid + '/envvariables',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    deleteVariable(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'envvariable/' + uuid + '/delete',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getVariable(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                PIPELINE_API_URL + 'envvariable/' + uuid,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updateVariable(uuid, name, value) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'envvariable/' + uuid,
                { name, value },
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    addVariable(name, value, organization_uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'envvariable',
                { name, value, organization_uuid },
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then((response) => {
                return response.data;
            });
    }
}

export default new EnvVariableService();
