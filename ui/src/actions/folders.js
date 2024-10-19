import {
    FOLDER_UPDATE_SUCCESS,
    FOLDER_UPDATE_FAIL,
    FOLDER_CREATE_FAIL,
    FOLDER_CREATE_SUCCESS,
    FOLDER_DELETE_SUCCESS,
    FOLDER_DELETE_FAIL,
    SET_MESSAGE,
} from "./types";

import {
    prepareErrorMessage,
} from "../actions/errors";

import FolderService from "../services/folder.service";

export const addFolder = (name, parentUuid, organizationUuid, type) => (dispatch) => {
    return FolderService.addFolder(name, parentUuid, organizationUuid, type).then(
        (data) => {
            dispatch({
                type: FOLDER_CREATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Folder successfully created",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: FOLDER_CREATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const updateFolder = (uuid, name, parentUuid) => (dispatch) => {
    return FolderService.updateFolder(uuid, name, parentUuid).then(
        (data) => {
            dispatch({
                type: FOLDER_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Folder successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: FOLDER_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const deleteFolder = (uuid) => (dispatch) => {
    return FolderService.deleteFolder(uuid).then(
        (data) => {
            dispatch({
                type: FOLDER_DELETE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Folder successfully deleted",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: FOLDER_DELETE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};
