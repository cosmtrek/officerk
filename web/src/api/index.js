import axios from 'axios';

export class API {
	constructor() {
		this.api = 'http://localhost:9392'
	}

	getEnum = () => {
		return axios.get(`${this.api}/enum`)
	}

	getNodesOnline = () => {
		return axios.get(`${this.api}/nodes/online`)
	}

	getNodeList = () => {
		return axios.get(`${this.api}/nodes`)
	}

	createNode = (params) => {
		return axios.post(`${this.api}/nodes`, params)
			.then((data) => data)
	}

	getNode = (id) => {
		return axios.get(`${this.api}/nodes/${id}`)
			.then((data) => data)
	}

	editNode = (id, params) => {
		return axios.put(`${this.api}/nodes/${id}`, params)
			.then((data) => data)
	}

	deleteNode = (id) => {
		return axios.delete(`${this.api}/nodes/${id}`)
			.then((data) => data)
	}

	getJobList = () => {
		return axios.get(`${this.api}/jobs`)
			.then((data) => data)
	}

	getJob = (id) => {
		return axios.get(`${this.api}/jobs/${id}`)
			.then((data) => data)
	}

	createJob = (params) => {
		return axios.post(`${this.api}/jobs`, params)
			.then((data) => data)
	}

	editJob = (id, params) => {
		return axios.put(`${this.api}/jobs/${id}`, params)
			.then((data) => data)
	}

	deleteJob = (id) => {
		return axios.delete(`${this.api}/jobs/${id}`)
			.then((data) => data)
	}

	getAJoblogs = (id) => {
		return axios.get(`${this.api}/jobs/${id}/logs`)
			.then((data) => data)
	}

	getJoblogList = () => {
		return axios.get(`${this.api}/joblogs`)
			.then((data) => data)
	}

	getJoblog = (id) => {
		return axios.get(`${this.api}/joblogs/${id}`)
			.then((data) => data)
	}
}