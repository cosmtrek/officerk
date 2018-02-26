import React from 'react'
import {
	NavLink as Link,
} from 'react-router-dom'
import {Row, Col, Switch, Modal, Button} from 'antd'
import {API} from "../../api"
import {G6Graph} from "../G6"
import '../../App.css'

export class Job extends React.Component {
	state = {
		api: new API(),
		id: '',
		job: null,
	}

	componentWillMount() {
		const id = this.props.match.params.id
		this.state.api.getJob(id)
			.then((resp) => {
				if (resp.error) {
					console.log(resp)
				} else {
					this.setState({
						id: id,
						job: resp.data.data,
					})
				}
			})
	}

	onDeleteNode = () => {
		const {id} = this.state
		const me = this
		Modal.confirm({
			title: "Do you want to delete this job?",
			onOk() {
				me.state.api.deleteJob(id)
					.then((resp) => {
						if (resp.error) {
							Notification.error('Delete Job', resp.error)
						} else {
							Notification.success('Delete Job', 'OK')
						}
					})
			},
			onCancel() {
				return me.props.history.goBack()
			},
		})
	}

	render() {
		const {job} = this.state
		if (job) {
			return (
				<div>
					<Row>
						<Col span={20}>
							<h2>Job #{job.id}@{job.node.ip ? job.node.ip : '[deleted node]'}: {job.name} {job.slug ? `(${job.slug})` : ''}</h2>
						</Col>
						<Col span={2} style={{textAlign: 'right'}}>
							<Link to={`/jobs/${job.id}/edit`}>
								<Button type="normal" icon="edit">Edit</Button>
							</Link>
						</Col>
						<Col span={2} style={{textAlign: 'right'}}>
							<Button
								type="danger"
								icon="delete"
								htmlType="submit"
								onClick={this.onDeleteNode}
							>
								Delete
							</Button>
						</Col>
					</Row>
					<Row className="dataContainer">
						<h3>Basic Info</h3>
						<Col span={3}>
							Status<span className="jobInfo"><Switch disabled checked={job.is_online}/></span>
						</Col>
						<Col span={3}>
							Type: <span className="jobInfo">{job.typ === 0 ? 'Cron' : 'Manual'}</span>
						</Col>
						{job.typ === 0 ? (
							<Col span={5}>
								Schedule: <span className="jobInfo">{job.schedule}</span>
							</Col>
						) : (
							<Col span={6}>
								Command: <span className="jobInfo">{`http://MASTER_IP/jobs/${job.id}/run`}</span>
							</Col>
						)}
						<Col span={5}>
							Created: <span className="jobInfo">{job.created_at}</span>
						</Col>
						<Col span={5}>
							Updated: <span className="jobInfo">{job.updated_at}</span>
						</Col>
					</Row>
					<Row className="dataContainer">
						<Col>
							<h3>Task List</h3>
							<ul className="taskList">
								{
									job.tasks.map((k, idx) => {
										return (
											<li
												className="jobInfo"
												key={idx}
											>
												<span className="taskInfo">
													{k.name}
												</span>
												<span className="taskInfo">
													{k.command}
												</span>
												<span className="taskInfo">
													{k.next_tasks}
												</span>
											</li>
										)
									})
								}
							</ul>
						</Col>
					</Row>
					<Row className="dataContainer">
						<Col>
							<h3>Tasks Dependency Graph</h3>
							<G6Graph
								nodes={this.drawNodes(job.graph.nodes)}
								edges={job.graph.edges}
							/>
						</Col>
					</Row>
				</div>
			)
		}
		return null
	}

	drawNodes = (nodes) => {
		nodes.forEach((n) => {
			n.id = n.name
			n.shape = 'rect'
			n.label = n.name
		})
		return nodes
	}
}
