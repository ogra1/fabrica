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

    repoList: (cancelCallback) => {
        return axios.get(constants.baseUrl + 'repos');
    },

    repoCreate: (repo, cancelCallback) => {
        return axios.post(constants.baseUrl + 'repos', {repo: repo});
    },
}

export default service