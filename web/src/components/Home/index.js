import React from 'react'
import {Row, Col, Card} from 'antd'

export class Home extends React.Component {
	render() {
		return (
			<div>
				<Row>
					<Col>
						<h2>Home</h2>
					</Col>
				</Row>
				<Row>
					<Col span={6}>
						<Card title="Node">
							<p>Card content</p>
						</Card>
					</Col>
					<Col span={6}>
						<Card title="Job">
							<p>Card content</p>
						</Card>
					</Col>
				</Row>
			</div>
		)
	}
}
