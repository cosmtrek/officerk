import React from 'react'
import {Form, Input, Button, Select, Switch, Icon} from 'antd'
import {TaskForm} from "../Task/form"
import {Notification} from "../Notification"

const FormItem = Form.Item
const Option = Select.Option
let taskId = 0

export class JobForm extends React.Component {
	handleSubmit = (e) => {
		const {op, job} = this.props
		e.preventDefault()
		this.props.form.validateFields((err, values) => {
			if (!err) {
				const params = values
				const tasks = []
				if (!('tasks' in params)) {
					// error
					return
				}
				Object.entries(params.tasks).forEach((v) => {
					tasks.push(v[1])
				})
				params.tasks = tasks
				delete params.task_keys
				console.log("params: ", params)
				console.log(JSON.stringify(params))
				if (op === 'CREATE') {
					this.props.api.createJob(params)
						.then((resp) => {
							if (resp.error) {
								Notification.error('Create Job', resp.error)
							} else {
								Notification.success('Create Job', 'OK')
							}
						})
				} else if (op === 'PUT') {
					params.id = job.id
					this.props.api.editJob(job.id, params)
						.then((resp) => {
							if (resp.error) {
								Notification.error('Edit Job', resp.error)
							} else {
								Notification.success('Edit Job', 'OK')
							}
						})
				}
			}
		})
	}

	removeTask = (k) => {
		const {form} = this.props
		const taskKeys = form.getFieldValue('task_keys')
		if (taskKeys.length === 1) {
			return
		}
		form.setFieldsValue({
			task_keys: taskKeys.filter(task => task !== k),
		})
	}

	addTask = () => {
		const {form} = this.props
		const taskKeys = form.getFieldValue('task_keys')
		const nextTasks = taskKeys.concat(`n${taskId}`)
		taskId++
		// can use data-binding to set
		// important! notify form to detect changes
		form.setFieldsValue({
			task_keys: nextTasks,
		})
	}

	render() {
		const {nodeList, job} = this.props
		const {getFieldDecorator, getFieldValue} = this.props.form
		const formItemLayout = {
			labelCol: { span: 2 },
			wrapperCol: { span: 6 },
		}
		const formTailLayout = {
			labelCol: { span: 2 },
			wrapperCol: { span: 6, offset: 2 },
		}
		const switchLayout = {
			colon: false,
			labelCol: {span: 2},
		}
		const taskLayout = {
			labelCol: {span: 2},
			wrapperCol: {span: 6},
		}
		const taskWithoutLabelLayout = {
			wrapperCol: {span: 6, offset: 2},
		}
		const addTaskLayout = {
			wrapperCol: {span: 10, offset: 2},
		}

		getFieldDecorator('task_keys', {initialValue: job ? job.tasks.map((k) => `o${k.id}`) : []})
		const taskKeys = getFieldValue('task_keys')
		const taskFormItems = taskKeys.map((k, index) => {
			return (
				<FormItem
					{...(index === 0 ? taskLayout : taskWithoutLabelLayout)}
					label={index === 0 ? 'Tasks' : ''}
					required={false}
					key={k}
				>
					{taskKeys.length > 1 ? (
						<Icon
							className="dynamic-delete-button"
							type="close-circle-o"
							disabled={taskKeys.length === 1}
							onClick={() => this.removeTask(k)}
						/>
					) : null}
					{getFieldDecorator(`tasks[${k}]`, {
						validateTrigger: ['onChange'],
						rules: [{
							required: true,
							message: "Please input task or delete it"
						}],
						initialValue: job ? job.tasks[index] : {},
					})(
						<TaskForm/>
					)}
				</FormItem>
			)
		})

		return (
			<Form onSubmit={this.handleSubmit}>
				<FormItem
					{...formItemLayout}
					label="Name"
				>
					{getFieldDecorator('name', {
						rules: [{required: true, message: 'Please input job name!'}],
						initialValue: job ? job.name : '',
					})(
						<Input/>
					)}
				</FormItem>
				<FormItem
					{...formItemLayout}
					label="Slug"
				>
					{getFieldDecorator('slug', {
						initialValue: job ? job.slug : '',
					})(
						<Input/>
					)}
				</FormItem>
				<FormItem
					{...formItemLayout}
					label="Schedule"
				>
					{getFieldDecorator('schedule', {
						initialValue: job ? job.schedule : '',
					})(
						<Input/>
					)}
				</FormItem>
				<FormItem
					{...formItemLayout}
					label="Type"
				>
					{getFieldDecorator('typ', {
						initialValue: job ? job.typ : 0,
					})(
						<Select>
							<Option key='cron' value={0}>Cron</Option>
							<Option key='manual' value={1}>Manual</Option>
						</Select>
					)}
				</FormItem>
				<FormItem
					{...switchLayout}
					label="Online"
				>
					{getFieldDecorator('is_online', {
						valuePropName: 'checked',
						initialValue: job ? job.is_online : false,
					})(
						<Switch/>
					)}
				</FormItem>
				<FormItem
					{...formItemLayout}
					label="Node IP"
				>
					{getFieldDecorator('node_id', {
						initialValue: job ? job.node_id : 'Node',
					})(
						<Select>
							{nodeList.map((node) => {
								return <Option key={node.id} value={node.id}>{node.ip}</Option>
							})}
						</Select>
					)}
				</FormItem>
				{taskFormItems}
				<FormItem
					{...addTaskLayout}
				>
					<Button type="dashed" onClick={this.addTask} style={{width: '60%'}}>
						<Icon type="plus"/> Add Task
					</Button>
				</FormItem>
				<FormItem
					{...formTailLayout}
				>
					<Button type="primary" htmlType="submit">
						{job ? 'Edit' : 'Create'}
					</Button>
				</FormItem>
			</Form>
		)
	}
}
