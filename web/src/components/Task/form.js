import React from 'react'
import {Input} from 'antd'

export class TaskForm extends React.Component {
	constructor(props) {
		super(props)

		const value = this.props.value || {}
		this.state = {
			name: value.name || '',
			command: value.command || '',
			nextTasks: value.nextTasks || value.next_tasks || '',
		}
	}

	componentWillReceiveProps(nextProps) {
		if ('value' in nextProps) {
			const value = nextProps.value
			this.setState(value)
		}
	}

	triggerChange = (changedValue) => {
		const onChange = this.props.onChange
		if (onChange) {
			onChange(Object.assign({}, this.state, changedValue))
		}
	}

	handleNameChange = (e) => {
		const name = e.target.value || ''
		if (!('value' in this.props)) {
			this.setState({name})
		}
		this.triggerChange({name})
	}


	handleCommandChange = (e) => {
		const command = e.target.value || ''
		if (!('value' in this.props)) {
			this.setState({command})
		}
		this.triggerChange({command})
	}

	handleNextTasksChange = (e) => {
		const nextTasks = e.target.value || ''
		if (!('value' in this.props)) {
			this.setState({nextTasks})
		}
		this.triggerChange({nextTasks})
	}

	render() {
		const state = this.state
		return (
			<div>
				<Input
					type="text"
					placeholder="Name"
					value={state.name}
					onChange={this.handleNameChange}
				/>
				<Input
					type="text"
					placeholder="Command"
					value={state.command}
					onChange={this.handleCommandChange}
				/>
				<Input
					type="text"
					placeholder="Next Tasks"
					value={state.nextTasks}
					onChange={this.handleNextTasksChange}
				/>
			</div>
		)
	}
}
