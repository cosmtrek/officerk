import React from 'react'
import {
	NavLink as Link,
} from 'react-router-dom'
import {Row, Col, Table, Icon} from 'antd'
import {API} from "../../api"
import "../../App.css"

const tableColumns = [
	{
		title: 'ID', dataIndex: 'id', key: 'id', render: (text, record) => {
			return <Link to={`/joblogs/${record.id}`}>{record.id}</Link>
		}
	},
	{title: 'Status', dataIndex: 'status', key: 'status', render: (text, record) => {
			return record.status === 1 ?
				<Icon type="check-circle" style={{fontSize: '16px', color: '#52c41a'}}/> :
				<Icon type="close-circle" style={{fontSize: '16px', color: '#f5222d'}}/>
		}},
	{title: 'Created At', dataIndex: 'created_at', key: 'created_at'},
	{title: 'Updated At', dataIndex: 'updated_at', key: 'updated_at'},
]

export class AJobLogs extends React.Component {
	state = {
		api: new API(),
		jobID: '',
		joblogList: [],
	}

	componentDidMount() {
		const id = this.props.match.params.id
		this.setState({
			jobID: id,
		})

		this.state.api.getAJoblogs(id)
			.then((resp) => {
				if (resp.error) {
					console.log(resp.error)
				} else {
					this.setState({
						joblogList: resp.data.data,
					})
				}
			})
	}

	render() {
		return (
			<div>
				<Row>
					<Col span={22}>
						<h2>Job Logs</h2>
					</Col>
				</Row>
				<Row className="dataContainer">
					<Col>
						<Table rowKey={record => record.id} columns={tableColumns} dataSource={this.state.joblogList}/>
					</Col>
				</Row>
			</div>
		)
	}
}
