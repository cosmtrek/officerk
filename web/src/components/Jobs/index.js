import React from 'react'
import {
	NavLink as Link,
} from 'react-router-dom'
import {Row, Col, Button, Table, Icon, Dropdown, Menu} from 'antd'
import {API} from "../../api"
import '../../App.css'

const tableColumns = [
	{title: 'ID', dataIndex: 'id', key: 'id'},
	{title: 'On', key: 'online', render: (text, record) => {
			return record.is_online ?
				<Icon type="check" style={{fontSize: '16px', color: '#52c41a'}}/> :
				<Icon type="close" style={{fontSize: '16px', color: '#f5222d'}}/>
		}},
	{title: 'Name', dataIndex: 'name', key: 'name', render: (text, record) => {
			return <Link to={`/jobs/${record.id}`}>{record.name}</Link>
		}},
	{
		title: 'Type', dataIndex: 'typ', key: 'typ', render: (text, record) => {
			return record.typ === 0 ? 'Cron' : 'Manual'
		}
	},
	{title: 'Slug', dataIndex: 'slug', key: 'slug'},
	{title: 'Schedule', dataIndex: 'schedule', key: 'schedule'},
	{title: 'Tasks', dataIndex: 'tasks', key: 'tasks', render: (text, record) => record.tasks.length},
	{
		title: 'Action', key: 'action', render: (text, record) => {
			const menu = (
				<Menu>
					{record.typ === 1 ? (
						<Menu.Item>
							<Link to={`/jobs/${record.id}/run`}>
								<Icon type="play-circle-o"/> Run
							</Link>
						</Menu.Item>
					) : null}
					<Menu.Item>
						<Link to={`/jobs/${record.id}/edit`}>
							<Icon type="edit"/> Edit
						</Link>
					</Menu.Item>
					<Menu.Item>
						<Link to={`/jobs/${record.id}/logs`}>
							<Icon type="dot-chart"/> Job Logs
						</Link>
					</Menu.Item>
				</Menu>
			)
			return (
				<div>
					<Dropdown overlay={menu}>
						<a className="ant-dropdown-link">
							Op <Icon type="down" />
						</a>
					</Dropdown>
				</div>
			)
		}
	},
]

export class Jobs extends React.Component {
	state = {
		api: new API(),
		jobList: [],
	}

	componentWillMount() {
		this.state.api.getJobList()
			.then((resp) => {
				if (resp.error) {
					console.log(resp.error)
				} else {
					this.setState({
						jobList: resp.data.data,
					})
				}
			})
	}

	render() {
		return (
			<div>
				<Row>
					<Col span={20}>
						<h2>Jobs</h2>
					</Col>
					<Col span={4}>
						<div style={{textAlign: 'right'}}>
							<Link to='/jobs/new'>
								<Button type="primary">Create Job</Button>
							</Link>
						</div>
					</Col>
				</Row>
				<Row className="dataContainer">
					<Col>
						<Table rowKey={record => record.id} columns={tableColumns} dataSource={this.state.jobList}/>
					</Col>
				</Row>
			</div>
		)
	}
}
