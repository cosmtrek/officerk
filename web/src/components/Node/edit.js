import React from 'react'
import {Row, Col, Form, Button, Modal} from 'antd'
import axios from 'axios'
import {NodeForm} from "./form"
import {API} from "../../api"
import {Notification} from "../Notification"

import "../../App.css"

export class EditNode extends React.Component {
	state = {
		api: new API(),
		op: 'PUT',
		nodeID: '',
		node: null,
		nodeIPs: [],
	}

	componentWillMount() {
		const id = this.props.match.params.id
		this.setState({
			nodeID: id,
		})

		axios.all([
			this.state.api.getNodesOnline(),
			this.state.api.getNode(id),
		]).then((arr) => {
			this.setState({
				nodeIPs: arr[0].data.data,
				node: arr[1].data.data,
			})
		})
	}

	onDeleteNode = () => {
		const {nodeID} = this.state
		const me = this
		Modal.confirm({
			title: "Do you want to delete this node?",
			onOk() {
				me.state.api.deleteNode(nodeID)
					.then((resp) => {
						if (resp.error) {
							Notification.error('Delete Node', resp.error)
						} else {
							Notification.success('Delete Node', 'OK')
						}
					})
			},
			onCancel() {
				return me.props.history.goBack()
			},
		})
	}

	render() {
		const EditNodeForm = Form.create({})(NodeForm)

		return (
			<div>
				<Row>
					<Col span={20}>
						<h2>Edit Node</h2>
					</Col>
					<Col span={4}>
						<div style={{textAlign: 'right'}}>
							<Button
								type="danger"
								icon="delete"
								htmlType="submit"
								onClick={this.onDeleteNode}
							>
								Delete
							</Button>
						</div>
					</Col>
				</Row>
				<Row>
					<Col className="formContainer">
						<EditNodeForm
							{...this.state}
						/>
					</Col>
				</Row>
			</div>
		)
	}
}
