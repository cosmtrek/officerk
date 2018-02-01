import {notification} from 'antd'

export class Notification {
	static success = (message, description) => {
		notification['success']({
			message: message,
			description: description,
		})
	}

	static error = (message, description) => {
		notification['error']({
			message: message,
			description: description,
		})
	}
}