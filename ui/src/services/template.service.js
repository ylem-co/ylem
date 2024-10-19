import axios from "axios";

import { PIPELINE_ROOT_FOLDER } from "./pipeline.service"; 

const TEMPLATE_API_URL_PREFIX = "/pipeline-api/";

const TEMPLATE_API_URL = TEMPLATE_API_URL_PREFIX + "pipeline-templates/";

export const TEMPLATES_LIST_TYPE_SYSTEM = "system";
export const TEMPLATES_LIST_TYPE_SHARED = "shared";
export const TEMPLATES_LIST_TYPE_ORG = "org";
export const TEMPLATES_LIST_TYPE_ME = "me";

class TemplateService {
	/*
		/pipeline-templates/ — shared templates of all organizations, TODO later
		/pipeline-templates/?onlySystem=1 — only templates created and shared by Ylem
		/organization/{uuid}/pipeline-templates/ — all templates of my organization
		/organization/{uuid}/pipeline-templates/?onlyMy=1 — only my templates
	*/
    getTemplates = async(type, orgUuid = null) => {
        var token = localStorage.getItem("token");
        let url = TEMPLATE_API_URL_PREFIX + "organization/" + orgUuid + "/pipeline-templates/";

        if (type === TEMPLATES_LIST_TYPE_ME) {
	        url = url + "?onlyMy=1";
        } else if (type === TEMPLATES_LIST_TYPE_SYSTEM) {
	        url = TEMPLATE_API_URL + "?onlySystem=1";
        }

        return axios
            .get(
                url,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    saveAsTemplate(pipeline_uuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                TEMPLATE_API_URL,
                { pipeline_uuid },
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then((response) => {
                return response.data;
            });
    }

    createFromTemplate(template_uuid, folder_uuid) {
        var token = localStorage.getItem("token");

        if (folder_uuid === PIPELINE_ROOT_FOLDER) {
        	folder_uuid = null;
        }

        return axios
            .post(
                TEMPLATE_API_URL + template_uuid + '/pipelines/',
                { folder_uuid },
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then((response) => {
                return response.data;
            });
    }
}

export default new TemplateService();
