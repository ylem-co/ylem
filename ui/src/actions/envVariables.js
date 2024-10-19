import {
    ENV_VARIABLE_UPDATE_SUCCESS,
    ENV_VARIABLE_UPDATE_FAIL,
    ENV_VARIABLE_CREATE_FAIL,
    ENV_VARIABLE_CREATE_SUCCESS,
    ENV_VARIABLE_DELETE_SUCCESS,
    ENV_VARIABLE_DELETE_FAIL,
    SET_MESSAGE,
} from "./types";

import {
    prepareErrorMessage,
} from "../actions/errors";

import EnvVariableService from "../services/envVariable.service";

export const addEnvVariable = (name, value, organizationUuid) => (dispatch) => {
    return EnvVariableService.addVariable(name, value, organizationUuid).then(
        (data) => {
            dispatch({
                type: ENV_VARIABLE_CREATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Environment variable successfully created",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: ENV_VARIABLE_CREATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const updateEnvVariable = (uuid, name, value) => (dispatch) => {
    return EnvVariableService.updateVariable(uuid, name, value).then(
        (data) => {
            dispatch({
                type: ENV_VARIABLE_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Environment variable successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: ENV_VARIABLE_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const deleteEnvVariable = (uuid) => (dispatch) => {
    return EnvVariableService.deleteVariable(uuid).then(
        (data) => {
            dispatch({
                type: ENV_VARIABLE_DELETE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Environment variable successfully deleted",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: ENV_VARIABLE_DELETE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};
