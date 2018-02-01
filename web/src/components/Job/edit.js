import React from 'react'
import {Row, Col, Form} from 'antd'
import axios from 'axios'
import {API} from "../../api"
import {JobForm} from "./form"
import '../../App.css'

export class EditJob extends React.Component {
	state = {
		api: new API(),
		op: 'PUT',
		jobID: '',
		nodeList: [],
	}

	componentDidMount() {
		const id = this.props.match.params.id
		this.setState({
			jobID: id,
		})

		axios.all([
			this.state.api.getNodeList(),
			this.state.api.getJob(id),
		]).then((arr) => {
				this.setState({
					nodeList: arr[0].data.data,
					job: arr[1].data.data,
				})
		})
	}

	render() {
		const EditJobForm = Form.create({})(JobForm)

		return (
			<div>
				<Row>
					<Col>
						<h2>Edit Job</h2>
					</Col>
				</Row>
				<Row>
					<Col className="formContainer">
						<EditJobForm
							{...this.state}
						/>
					</Col>
				</Row>
			</div>
		)
	}
}
