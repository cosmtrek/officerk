import React from 'react'
import PropTypes from 'prop-types'
import G6 from '@antv/g6'
import Plugins from '@antv/g6-plugins'

G6.track(false)

function noop() {}

const dagre = new Plugins['layout.dagre']({
	rankdir: 'LR',
	nodesep: 120,
	ranksep: 120,
})

// Copy from https://codesandbox.io/s/p3jry6w230
export class G6Graph extends React.Component {

	static defaultProps = {
		plugins: [dagre],
		className: 'g6-graph',
		nodeLabelFill: '#0B73B3',
		nodeShapeStroke: '#00B7ED',
		groupLabelFill: '#DA7639',
		groupShapeStroke: '#FABE9C',
		edgeStroke: '#00AEEC',
		forceFit: true,
		height: 500,
		grid: {
			forceAlign: true,
			cell: 10
		},
		onClick: noop,
		onMouseDown: noop,
		onMouseMove: noop,
		onMouseUp: noop,
		onMouseLeave: noop,
		onMouseEnter: noop,
	}

	static propTypes = {
		nodes: PropTypes.array.isRequired,
		edges: PropTypes.array.isRequired,
		plugins: PropTypes.array,
		className: PropTypes.string,
		nodeLabelFill: PropTypes.string,
		nodeShapeStroke: PropTypes.string,
		groupLabelFill: PropTypes.string,
		groupShapeStroke: PropTypes.string,
		edgeStroke: PropTypes.string,
		forceFit: PropTypes.bool,
		height: PropTypes.number,
		grid: PropTypes.object,
		onClick: PropTypes.func,
		onMouseDown: PropTypes.func,
		onMouseMove: PropTypes.func,
		onMouseUp: PropTypes.func,
		onMouseLeave: PropTypes.func,
		onMouseEnter: PropTypes.func,
	}

	_renderG6Graph() {
		const me = this
		const container = me.container
		const props = me.props

		if (me.graph) {
			me.graph.destroy()
		}

		const net = new G6.Net({
			container,
			...props,
		})

		net.source(props.nodes, props.edges)
		net.edge().shape('arrow')

		net.on('click', props.onClick.bind(me))
		net.on('mousedown', props.onMouseDown.bind(me))
		net.on('mousemove', props.onMouseMove.bind(me))
		net.on('mouseup', props.onMouseUp.bind(me))
		net.on('mouseenter', props.onMouseEnter.bind(me))
		net.on('mouseleave', props.onMouseLeave.bind(me))

		net.render()

		me.net = me.g6graph = net
	}

	componentDidMount() {
		this._renderG6Graph()
	}

	componentDidUpdate() {
		this._renderG6Graph()
	}

	render() {
		const me = this
		const props = me.props
		return (
			<div
				className={props.className}
				ref={(container) => {
					me.container = container
				}}>
			</div>
		)
	}
}
