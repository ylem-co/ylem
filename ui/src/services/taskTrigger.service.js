import axios from "axios";

const PIPELINE_API_URL = "/pipeline-api/";

export const TRIGGER_TYPE_SCHEDULE = "schedule"
export const TRIGGER_TYPE_CONDITION_TRUE = "condition_true"
export const TRIGGER_TYPE_CONDITION_FALSE = "condition_false"
export const TRIGGER_TYPE_OUTPUT = "output"

class TaskTriggerService {
    getTaskTriggersByPipeline = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                PIPELINE_API_URL + 'pipeline/' + uuid + '/task-triggers',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    deleteTaskTrigger(uuid, pipelineUuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + pipelineUuid + '/task-trigger/' + uuid + '/delete',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getTaskTrigger(uuid, pipelineUuid) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                PIPELINE_API_URL + 'pipeline/' + pipelineUuid + '/task-trigger/' + uuid,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updateTaskTrigger(uuid, pipelineUuid, trigger_type, schedule = "") {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + pipelineUuid + '/task-trigger/' + uuid,
                { trigger_type, schedule },
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    addTaskTrigger(pipelineUuid, trigger_task_uuid, triggered_task_uuid, trigger_type, schedule = "") {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + pipelineUuid + '/task-trigger',
                { trigger_task_uuid, triggered_task_uuid, trigger_type, schedule },
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then((response) => {
                return response.data;
            });
    }
}

export default new TaskTriggerService();
