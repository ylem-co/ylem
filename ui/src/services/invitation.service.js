import axios from "axios";

const API_URL = "/user-api/";

class InvitationService {
    getPendingByOrganization = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                API_URL + 'organization/' + uuid + '/pending-invitations',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    sendInvitations = async(uuid, emails) => {
        var token = localStorage.getItem("token");

        return axios
            .post(
                API_URL + 'organization/' + uuid + '/invitations',
                {emails},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    validateInvitationKey = async(key) => {
        return axios
            .post(
                API_URL + 'invitations/' + key + '/validate',
                {},
                {}
            );
    }
}

export default new InvitationService();
