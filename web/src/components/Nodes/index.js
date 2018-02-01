import React from 'react'
import {
	NavLink as Link,
} from 'react-router-dom'
import {Row, Col, Button, Table, Icon, Dropdown, Menu} from 'antd'
import {API} from "../../api"
import '../../App.css'

const tableColumns = [
	{title: 'ID', dataIndex: 'id', key: 'id'},
	{title: 'Online', key: 'online', render: (text, record) => {
		return record.online ?
			<Icon type="check" style={{fontSize: '16px', color: '#52c41a'}}/> :
			<Icon type="close" style={{fontSize: '16px', color: '#f5222d'}}/>
		}},
	{title: 'Name', dataIndex: 'name', key: 'name'},
	{title: 'IP', dataIndex: 'ip', key: 'ip'},
	{title: 'Action', key: 'action', render: (text, record) => {
		const menu = (
				<Menu>
					<Menu.Item>
						<Link to={`/nodes/${record.id}/edit`}>
							<Icon type="edit"/> Edit
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
		}},
]

export class Nodes extends React.Component {
	state = {
		api: new API(),
		nodeList: [],
	}

	componentWillMount() {
		this.state.api.getNodeList()
			.then((resp) => {
				if (resp.error) {
					console.log(resp.error)
				}	else {
					this.setState({
						nodeList: resp.data.data,
					})
				}
		})
	}

	render() {
		return (
			<div>
				<Row>
					<Col span={20}>
						<h2>Nodes</h2>
					</Col>
					<Col span={4}>
						<div style={{textAlign: 'right'}}>
							<Link to='/nodes/new'>
								<Button type="primary">Create Node</Button>
							</Link>
						</div>
					</Col>
				</Row>
				<Row className="dataContainer">
					<Col>
						<Table rowKey={record => record.id} columns={tableColumns} dataSource={this.state.nodeList} />
					</Col>
				</Row>
			</div>
		)
	}
}