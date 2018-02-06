import React from 'react'
import {Form, Input, Button, Select} from 'antd'
import {Notification} from "../Notification"

const FormItem = Form.Item
const Option = Select.Option

export class NodeForm extends React.Component {
	handleSubmit = (e) => {
		const {op, node} = this.props
		e.preventDefault()
		this.props.form.validateFields((err, values) => {
			if (!err) {
				console.log('params: ', values)
				const params = values
				if (op === 'CREATE') {
					this.props.api.createNode(params)
						.then((resp) => {
							Notification.success('Create Node', 'OK')
						})
						.catch((error) => {
							Notification.error('Create Node', error.message + error.response.error)
						})
				} else if (op === 'PUT') {
					this.props.api.editNode(node.id, params)
						.then((resp) => {
							Notification.success('Edit Node', 'OK')
						})
						.catch((error) => {
							Notification.error('Create Node', error.message + error.response.error)
						})
				}
			}
		})
	}

	render() {
		const {nodeIPs, node} = this.props
		const {getFieldDecorator} = this.props.form
		const formItemLayout = {
			labelCol: { span: 2 },
			wrapperCol: { span: 6 },
		}
		const formTailLayout = {
			labelCol: { span: 2 },
			wrapperCol: { span: 6, offset: 2 },
		}

		return (
			<Form onSubmit={this.handleSubmit}>
				<FormItem
					{...formItemLayout}
					label="Name"
				>
					{getFieldDecorator('name', {
						rules: [{required: true, message: 'Please input node name!'}],
						initialValue: node ? node.name : '',
					})(
						<Input/>
					)}
				</FormItem>
				<FormItem
					{...formItemLayout}
					label="Node IP"
				>
					{getFieldDecorator('ip', {
						rules: [{required: true, messsage: 'Please select a node'}],
						initialValue: node ? node.ip : '',
					})(
						<Select>
							{nodeIPs.map((ip) => {
								return <Option key={ip} value={ip}>{ip}</Option>
							})}
						</Select>
					)}
				</FormItem>
				<FormItem
					{...formTailLayout}
				>
					<Button type="primary" htmlType="submit">
						{node ? 'Edit' : 'Create'}
					</Button>
				</FormItem>
			</Form>
		)
	}
}
