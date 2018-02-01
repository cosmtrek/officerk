import React from 'react'
import {Switch, Route} from 'react-router-dom'

import {AppLayout} from "../components/Layout"
import {Home} from "../components/Home"
import {Nodes} from "../components/Nodes"
import {CreateNode} from "../components/Node/create"
import {EditNode} from "../components/Node/edit"
import {Jobs} from "../components/Jobs"
import {CreateJob} from "../components/Job/create"
import {Job} from "../components/Job"
import {EditJob} from "../components/Job/edit"
import {AJobLogs} from "../components/Job/logs"
import {JobLogs} from "../components/Joblogs"
import {JobLog} from "../components/Joblog"


export class Routes extends React.Component {
	render() {
		return (
			<AppLayout>
				<Switch>
					<Route exact path='/' component={Home}/>
					<Route exact path='/nodes' component={Nodes}/>
					<Route exact path='/nodes/new' component={CreateNode}/>
					<Route exact path='/nodes/:id/edit' component={EditNode}/>
					<Route exact path='/jobs' component={Jobs}/>
					<Route exact path='/jobs/new' component={CreateJob}/>
					<Route exact path='/jobs/:id' component={Job}/>
					<Route exact path='/jobs/:id/edit' component={EditJob}/>
					<Route exact path='/jobs/:id/logs' component={AJobLogs}/>
					<Route exact path='/joblogs' component={JobLogs}/>
					<Route exact path='/joblogs/:id' component={JobLog}/>
				</Switch>
			</AppLayout>
		)
	}
}