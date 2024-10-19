import {
    TASK_TRIGGER_UPDATE_SUCCESS,
    TASK_TRIGGER_UPDATE_FAIL,
    TASK_TRIGGER_CREATE_FAIL,
    TASK_TRIGGER_DELETE_SUCCESS,
    TASK_TRIGGER_DELETE_FAIL,
    SET_MESSAGE,
} from "./types";

import {
    prepareErrorMessage,
} from "../actions/errors";

import TaskTriggerService from "../services/taskTrigger.service";

export const addTaskTrigger = (pipelineUuid, triggerTaskUuid, triggeredTaskUuid, triggerType, schedule = "") => (dispatch) => {
    return TaskTriggerService.addTaskTrigger(pipelineUuid, triggerTaskUuid, triggeredTaskUuid, triggerType, schedule).then(
        (data) => {
            dispatch({
                type: SET_MESSAGE,
                payload: { item: data },
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: TASK_TRIGGER_CREATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const updateTaskTrigger = (uuid, pipelineUuid, triggerType, schedule = "") => (dispatch) => {
    return TaskTriggerService.updateTaskTrigger(uuid, pipelineUuid, triggerType, schedule).then(
        (data) => {
            dispatch({
                type: TASK_TRIGGER_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Task trigger successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: TASK_TRIGGER_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const deleteTaskTrigger = (uuid, pipelineUuid) => (dispatch) => {
    return TaskTriggerService.deleteTaskTrigger(uuid, pipelineUuid).then(
        (data) => {
            dispatch({
                type: TASK_TRIGGER_DELETE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Task trigger successfully deleted",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: TASK_TRIGGER_DELETE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};
