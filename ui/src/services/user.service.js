import axios from "axios";

const API_URL = "/user-api/";

class UserService {
    updatePassword(uuid, password, confirm_password) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                API_URL + '/me/password' ,
                {password, confirm_password},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updateMe(uuid, first_name, last_name, email, phone) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                API_URL + 'me',
                {first_name, last_name, email, phone},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    assignRoleToUser(uuid, role) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                API_URL + 'user/' + uuid + "/assign-role",
                {role},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    terminateUser(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                API_URL + 'user/' + uuid + '/delete',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    activateUser(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                API_URL + 'user/' + uuid + '/activate',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    confirmEmail(key) {
        return axios
            .post(
                API_URL + 'email/' + key + '/confirm',
                {},
                {}
            );
    }

    getUsersByOrganization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                API_URL + 'organization/' + uuid + '/users',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }
}

export default new UserService();
