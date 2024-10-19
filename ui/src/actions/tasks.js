import {
    TASK_UPDATE_SUCCESS,
    TASK_UPDATE_FAIL,
    TASK_CREATE_FAIL,
    TASK_DELETE_SUCCESS,
    TASK_DELETE_FAIL,
    SET_MESSAGE,
} from "./types";

import {
    prepareErrorMessage,
} from "../actions/errors";

import TaskService from "../services/task.service";

export const addTask = (pipelineUuid, name, type) => (dispatch) => {
    return TaskService.addTask(pipelineUuid, name, type).then(
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
                type: TASK_CREATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const updateTask = (uuid, pipelineUuid, name, severity, type, implementation = null) => (dispatch) => {
    return TaskService.updateTask(uuid, pipelineUuid, name, severity, type, implementation).then(
        (data) => {
            dispatch({
                type: TASK_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Task successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: TASK_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const deleteTask = (uuid, pipelineUuid) => (dispatch) => {
    return TaskService.deleteTask(uuid, pipelineUuid).then(
        (data) => {
            dispatch({
                type: TASK_DELETE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Task successfully deleted",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: TASK_DELETE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};
