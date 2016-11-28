import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';

class App extends Component {
  constructor(props) {
    super(props)
    this.state = {
      messages: []
    }
  }
  componentDidMount() {
    this.connection = new WebSocket('ws://localhost:8083/ws');
    this.connection.onmessage = evt => { 
      // add the new message to state
        this.setState({
          messages : this.state.messages.concat([ evt.data ])
        })
    };

    // for testing: sending a message to the echo service every 2 seconds, 
    // the service sends it right back
    // setInterval( _ =>{
    //     this.connection.send( Math.random() )
    // }, 2000 )
  }
  onSubmit(e) {
    e.preventDefault()
    console.log(e.target.value)
    this.connection.send(this.state.message)
    this.setState({
      message: ''
    })
  }
  onChange(e) {
    this.setState({
      message: e.target.value
    })
  }
  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Welcome to React</h2>
        </div>
          <form onSubmit={this.onSubmit.bind(this)}>
            <ul>{ this.state.messages.map( (msg, idx) => <li key={'msg-' + idx }>{ msg }</li> )}</ul>;
            <input type="text" onChange={this.onChange.bind(this)} value={this.state.message} />
          </form>
      </div>
    );
  }
}

export default App;
