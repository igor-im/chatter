/**
 * Sample React Native App
 * https://github.com/facebook/react-native
 * @flow
 */

import React, { Component } from 'react';
import {
  AppRegistry,
  StyleSheet,
  Text,
  View,
  ListView,
  TextInput,
  Button
} from 'react-native';

export default class chatterandroid extends Component {
  constructor(props) {
    super(props)
    this.state = {
      messages: []
    }
  }
  componentDidMount() {
    this.connection = new WebSocket('ws://10.0.2.2:8083/ws');
    this.connection.onmessage = evt => { 
      // add the new message to state
        this.setState({
          messages : this.state.messages.concat([evt.data])
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
    this.connection.send(this.state.message)
    this.setState({
      message: ''
    })
  }
  onChange(text) {
    this.setState({
      message: text
    })
  }
  render() {
    return (
      <View style={styles.container}>
        <Text style={styles.welcome}>
          Welcome to React Native!
        </Text>
        <Text style={styles.instructions}>
          To get started, edit index.android.js
        </Text>
          {this.state.messages.map((v, i) => {
            return <Text>{v}</Text>
          })}
          <TextInput
          onChangeText={this.onChange.bind(this)} 
          value={this.state.message} />
          <Button
            onPress={this.onSubmit.bind(this)}
            title="submit"
            color="#841584"
            accessibilityLabel="Learn more about this purple button"
          />
      </View>
    );
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#F5FCFF',
  },
  welcome: {
    fontSize: 20,
    textAlign: 'center',
    margin: 10,
  },
  instructions: {
    textAlign: 'center',
    color: '#333333',
    marginBottom: 5,
  },
});

AppRegistry.registerComponent('chatterandroid', () => chatterandroid);
