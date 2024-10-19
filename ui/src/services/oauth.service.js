import axios from "axios";

const API_URL = "/oauth-api/";

class OAuthService {
    getClients = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                API_URL + 'clients/',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    createClient = async(name) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                API_URL + 'clients/',
                { name },
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    deleteClient = async(uuid) => {
        var token = localStorage.getItem("token")
        return axios
            .post(
                API_URL + 'client/' + uuid + '/delete/',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }
}

export default new OAuthService();
