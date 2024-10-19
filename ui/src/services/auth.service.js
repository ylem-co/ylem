import axios from "axios";

const API_URL = "/user-api/";

class AuthService {
    login(email, password) {
        return axios
            .post(API_URL + "login", { email, password })
            .then((response) => {
                if (response.data.token) {
                    localStorage.setItem("token", response.data.token);
                    const dataToSave = Object.assign({}, response.data);
                    delete dataToSave["token"];
                    localStorage.setItem("user", JSON.stringify(dataToSave));
                }

                return response.data;
            })
            .then((data) => {
                return axios.get(API_URL + 'my-organization', { headers: { Authorization: 'Bearer ' + data.token } });
            })
            .then((organization) => {
                if (organization.data && organization.data.name) {
                    localStorage.setItem("organization", JSON.stringify(organization.data));
                }

                return JSON.parse(localStorage.getItem("user"));
            });
    }

    logout = async() => {
        var token = localStorage.getItem("token");

        return axios.post(
            API_URL + "logout", 
            {},
            { headers: { Authorization: 'Bearer ' + token } }
        ).then((response) => {
            localStorage.removeItem("token");
            localStorage.removeItem("user");
            localStorage.removeItem("organization");
        }).catch(function (error) {
            if (
                error.response
                && error.response.status === 403
            ) {
                localStorage.removeItem("token");
                localStorage.removeItem("user");
                localStorage.removeItem("organization");
            }
        });
    }

    register(
        firstName,
        lastName,
        email, 
        password, 
        confirmPassword,
        phone, 
        organizationName = null,
        invitationKey = null
    ) {
        var data = {
            "first_name": firstName,
            "last_name": lastName,
            "email": email, 
            "password": password, 
            "confirm_password": confirmPassword,
            "phone": phone,
        };

        if (invitationKey !== null) {
            data.invitation_key = invitationKey;
        }

        if (organizationName !== null) {
            data.organization_name = organizationName;
        }

        return axios.post(API_URL + "user", data);
    }

    getSignInWithGoogleRedirectUrl = async() => {
        return axios
            .get(
                API_URL + 'auth/google',
                {
                    withCredentials: true
                }
            );
    }

    getSignInWithGoogleAvailability = async() => {
        return axios
            .get(
                API_URL + 'auth/google/available',
                {
                    withCredentials: true
                }
            );
    }

    signInWithGoogle = async(queryString) => {
        return axios
            .post(
                API_URL + 'auth/google/callback' + queryString,
                {},
                {
                    withCredentials: true
                }
            )
            .then((response) => {
                if (response.data.token) {
                    localStorage.setItem("token", response.data.token);
                    const dataToSave = Object.assign({}, response.data);
                    delete dataToSave["token"];
                    localStorage.setItem("user", JSON.stringify(dataToSave));
                }

                return response.data;
            })
            .then((data) => {
                return axios.get(API_URL + 'my-organization', { headers: { Authorization: 'Bearer ' + data.token } });
            })
            .then((organization) => {
                if (organization.data && organization.data.name) {
                    localStorage.setItem("organization", JSON.stringify(organization.data));
                }

                return JSON.parse(localStorage.getItem("user"));
            });
    }
}

export default new AuthService();
