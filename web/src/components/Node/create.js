import React from 'react'
import {Row, Col, Form} from 'antd'
import {API} from "../../api"
import {NodeForm} from "./form"

import "../../App.css"

export class CreateNode extends React.Component {
	state = {
		api: new API(),
		op: 'CREATE',
		nodeIPs: [],
	}

	componentDidMount() {
		this.state.api.getNodesOnline()
			.then((resp) => {
				if (resp.error) {
					console.log(resp.error)
				} else {
					this.setState({
						nodeIPs: resp.data.data,
					})
				}
			})
	}

	render() {
		const CreateNodeForm = Form.create({})(NodeForm)

		return (
			<div>
				<Row>
					<Col>
						<h2>Create Node</h2>
					</Col>
				</Row>
				<Row>
					<Col className="formContainer">
						<CreateNodeForm
							{...this.state}
						/>
					</Col>
				</Row>
			</div>
		)
	}
}
