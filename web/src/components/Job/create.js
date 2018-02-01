import React from 'react'
import {Row, Col, Form} from 'antd'
import {API} from "../../api"
import {JobForm} from "./form"
import '../../App.css'

export class CreateJob extends React.Component {
	state = {
		api: new API(),
		op: 'CREATE',
		nodeList: [],
	}

	componentDidMount() {
		this.state.api.getNodeList()
			.then((resp) => {
				if (resp.error) {
					console.log(resp.error)
				} else {
					this.setState({
						nodeList: resp.data.data,
					})
				}
			})
	}

	render() {
		const CreateJobForm = Form.create({})(JobForm)

		return (
			<div>
				<Row>
					<Col>
						<h2>Create Job</h2>
					</Col>
				</Row>
				<Row>
					<Col className="formContainer">
						<CreateJobForm
							{...this.state}
						/>
					</Col>
				</Row>
			</div>
		)
	}
}
