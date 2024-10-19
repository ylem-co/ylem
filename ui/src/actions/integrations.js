import {
    INTEGRATION_UPDATE_SUCCESS,
    INTEGRATION_UPDATE_FAIL,
    INTEGRATION_CREATE_FAIL,
    INTEGRATION_TEST_SUCCESS,
    INTEGRATION_TEST_FAIL,
    SET_MESSAGE,
    SLACK_AUTHORIZATION_UPDATE_FAIL,
    SLACK_AUTHORIZATION_UPDATE_SUCCESS,
    JIRA_AUTHORIZATION_UPDATE_FAIL,
    JIRA_AUTHORIZATION_UPDATE_SUCCESS,
    HUBSPOT_AUTHORIZATION_UPDATE_SUCCESS,
    HUBSPOT_AUTHORIZATION_UPDATE_FAIL, SALESFORCE_AUTHORIZATION_UPDATE_SUCCESS, SALESFORCE_AUTHORIZATION_UPDATE_FAIL,
} from "./types";

import {
    prepareErrorMessage,
} from "../actions/errors";

import IntegrationService from "../services/integration.service";
import log from "loglevel";

export const addIntegration = (type, data, sqlType = false) => (dispatch) => {
    return IntegrationService.addIntegration(type, data, sqlType).then(
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
                type: INTEGRATION_CREATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const updateIntegration = (uuid, type, data, sqlType = false) => (dispatch) => {
    return IntegrationService.updateIntegration(uuid, type, data, sqlType).then(
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
                type: INTEGRATION_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const confirmIntegration = (type, code, uuid) => (dispatch) => {
    return IntegrationService.confirmIntegration(type, code, uuid).then(
        (data) => {
            dispatch({
                type: INTEGRATION_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Integration successfully confirmed",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: INTEGRATION_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const updateSlackAuthorization = (uuid, name) => (dispatch) => {
    log.debug("update slack authorization")

    return IntegrationService.updateSlackAuthorization(name, uuid).then(
        (data) => {
            dispatch({
                type: SLACK_AUTHORIZATION_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Slack authorization successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: SLACK_AUTHORIZATION_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
}

export const updateSalesforceAuthorization = (uuid, name) => (dispatch) => {
    log.debug("update salesforce authorization")

    return IntegrationService.updateSalesforceAuthorization(name, uuid).then(
        (data) => {
            dispatch({
                type: SALESFORCE_AUTHORIZATION_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Salesforce authorization successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: SALESFORCE_AUTHORIZATION_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
}

export const updateJiraAuthorization = (uuid, name, resourceId) => (dispatch) => {
    log.debug("update jira authorization")

    return IntegrationService.updateJiraAuthorization(name, resourceId, uuid).then(
        (data) => {
            dispatch({
                type: JIRA_AUTHORIZATION_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Jira authorization successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: JIRA_AUTHORIZATION_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
}

export const updateHubspotAuthorization = (uuid, name) => (dispatch) => {
    log.debug("update hubspot authorization")

    return IntegrationService.updateHubspotAuthorization(name, uuid).then(
        (data) => {
            dispatch({
                type: HUBSPOT_AUTHORIZATION_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Hubspot authorization successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: HUBSPOT_AUTHORIZATION_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
}

export const testNewIntegration = (type, data) => (dispatch) => {
    return IntegrationService.testNewIntegration(type, data).then(
        (data) => {
            dispatch({
                type: INTEGRATION_TEST_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Connection succeeded",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);


            dispatch({
                type: INTEGRATION_TEST_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const testExistingIntegration = (uuid, type, data) => (dispatch) => {
    return IntegrationService.testExistingIntegration(uuid, type, data).then(
        (data) => {
            dispatch({
                type: INTEGRATION_TEST_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Connection succeeded",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: INTEGRATION_TEST_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};
