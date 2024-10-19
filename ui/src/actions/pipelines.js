import {
    PIPELINE_UPDATE_SUCCESS,
    PIPELINE_UPDATE_FAIL,
    PIPELINE_CREATE_FAIL,
    SET_MESSAGE,
} from "./types";

import {
    prepareErrorMessage,
} from "../actions/errors";

import PipelineService, { PIPELINE_TYPE_GENERIC } from "../services/pipeline.service";

export const addPipeline = (name, folderUuid, organizationUuid, layoutElements = [], type = PIPELINE_TYPE_GENERIC, schedule = null) => (dispatch) => {
    return PipelineService.addPipeline(name, folderUuid, organizationUuid, layoutElements, type, schedule).then(
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
                type: PIPELINE_CREATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const updatePipeline = (uuid, name, folderUuid, layoutElements = [], schedule = null) => (dispatch) => {
    return PipelineService.updatePipeline(uuid, name, folderUuid, layoutElements, schedule).then(
        (data) => {
            dispatch({
                type: PIPELINE_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Pipeline successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: PIPELINE_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};
