import React from 'react'
import {
	BrowserRouter as Router,
} from 'react-router-dom'

import {Routes} from "./routes"


class App extends React.Component {
	render() {
		return (
			<Router>
				<div>
					<Routes/>
				</div>
			</Router>
		)
	}
}

export default App
