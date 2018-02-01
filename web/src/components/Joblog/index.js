import React from 'react'
import {Row, Col, Timeline} from 'antd'
import {API} from "../../api"
import "../../App.css"

export class JobLog extends React.Component {
	state = {
		api: new API(),
		id: '',
		joblog: null,
	}

	componentWillMount() {
		const id = this.props.match.params.id
		this.state.api.getJoblog(id)
			.then((resp) => {
				if (resp.error) {
					console.log(resp)
				} else {
					this.setState({
						id: id,
						joblog: resp.data.data,
					})
				}
			})
	}

	render() {
		const {joblog} = this.state
		if (joblog) {
			return (
				<div>
					<Row>
						<Col>
							<h2>Job Log #{joblog.id} - {joblog.job.name}</h2>
						</Col>
					</Row>
					<Row className="dataContainer">
						<Col>
							<Timeline>
								{
									joblog.task_logs.map((k, idx) => {
										return (
											<Timeline.Item
												color={k.status === 1 ? 'green' : 'red'}
											>
												<p>{k.created_at}: {k.task.name}</p>
												<p>{k.result}</p>
											</Timeline.Item>
										)
									})
								}
							</Timeline>
						</Col>
					</Row>
				</div>
			)
		}
		return null
	}
}
