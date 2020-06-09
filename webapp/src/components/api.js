import axios from 'axios'
import constants from './constants'

let service = {
    build: (repo, cancelCallback) => {
        return axios.post(constants.baseUrl + 'build', {repo: repo});
    },

    buildList: (cancelCallback) => {
        return axios.get(constants.baseUrl + 'builds');
    },

    buildGet: (buildId, cancelCallback) => {
        return axios.get(constants.baseUrl + 'builds/' + buildId);
    },

    buildDelete: (buildId, cancelCallback) => {
        return axios.delete(constants.baseUrl + 'builds/' + buildId);
    },

    repoList: (cancelCallback) => {
        return axios.get(constants.baseUrl + 'repos');
    },

    repoCreate: (repo, branch, keyId, cancelCallback) => {
        return axios.post(constants.baseUrl + 'repos', {repo: repo, branch: branch, keyId: keyId});
    },

    repoDelete: (repoId, deleteBuilds, cancelCallback) => {
        return axios.post(constants.baseUrl + 'repos/delete', {id: repoId, deleteBuilds: deleteBuilds});
    },

    imageList: (cancelCallback) => {
        return axios.get(constants.baseUrl + 'check/images');
    },

    connectionList: (cancelCallback) => {
        return axios.get(constants.baseUrl + 'check/connections');
    },

    systemMonitor: (cancelCallback) => {
        return axios.get(constants.baseUrl + 'system');
    },

    systemEnvironment: (cancelCallback) => {
        return axios.get(constants.baseUrl + 'system/environment');
    },

    keysList: (cancelCallback) => {
        return axios.get(constants.baseUrl + 'keys');
    },

    keysCreate: (key, cancelCallback) => {
        return axios.post(constants.baseUrl + 'keys', key);
    },
}

export default service