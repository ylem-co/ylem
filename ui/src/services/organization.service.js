import axios from "axios";

const API_URL = "/user-api/";

class OrganizationService {
    updateOrganization(uuid, name) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                API_URL + 'organization/' + uuid,
                {name},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getMyOrganization = async() => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                API_URL + 'my-organization',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }
}

export default new OrganizationService();
