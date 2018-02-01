import React from 'react'
import {Layout, Menu, Icon} from 'antd'
import {
	NavLink as Link,
} from 'react-router-dom'

const {Header, Content, Footer, Sider} = Layout

export class AppLayout extends React.Component {
	state = {
		collapsed: false,
	}

	onCollapse = (collapsed) => {
		this.setState({collapsed})
	}

	render() {
		return (
			<Layout style={{minHeight: '100vh'}}>
				<Sider
					collapsible
					collapsed={this.state.collapsed}
					onCollapse={this.onCollapse}
					collapsedWidth={82}
					width={180}
				>
					<div id="logo" style={{textAlign: 'center', margin: '20px 0 20px 0'}}>
						<h2 style={{color: '#fff'}}>OfficerK</h2>
					</div>
					<Menu theme="dark" mode="inline">
						<Menu.Item key="home">
							<Icon type="home"/>
							<span><Link
								to='/'
								exact
								style={styles.menuItemStyle}
								activeStyle={styles.menuItemActiveStyle}
							>Home</Link></span>
						</Menu.Item>
						<Menu.Item key="node">
							<Icon type="cloud-o"/>
							<span><Link
								to='/nodes'
								style={styles.menuItemStyle}
								activeStyle={styles.menuItemActiveStyle}
							>Nodes</Link></span>
						</Menu.Item>
						<Menu.Item key="job">
							<Icon type="schedule"/>
							<span><Link
								to='/jobs'
								style={styles.menuItemStyle}
								activeStyle={styles.menuItemActiveStyle}
							>Jobs</Link></span>
						</Menu.Item>
						<Menu.Item key="joblog">
							<Icon type="dot-chart"/>
							<span><Link
								to='/joblogs'
								style={styles.menuItemStyle}
								activeStyle={styles.menuItemActiveStyle}
							>Job Logs</Link></span>
						</Menu.Item>
					</Menu>
				</Sider>
				<Layout>
					<Header style={{background: '#fff', padding: 0}}/>
					<Content style={{margin: '30px 20px'}}>
						<div style={{padding: 26, background: '#fff', minHeight: 360, textAlign: 'left'}}>
							{this.props.children}
						</div>
					</Content>
					<Footer style={{textAlign: 'center'}}>
						OfficerK © 2018 Made by Rick Yu
					</Footer>
				</Layout>
			</Layout>
		)
	}
}

const styles = {
	menuItemStyle: {
		color: '#999',
		marginLeft: '12px',
	},
	menuItemActiveStyle: {
		color: '#fff',
		textDecoration: 'none',
	}
}
