/** @jsx React.DOM */
var Todo = React.createClass({
    render: function() {
        return (
                <div className="todo">
                <h2 className="status">
                {this.props.status}
            </h2>
                {this.props.children}
            </div>
        );
    }
});
var TodoList = React.createClass({
    render: function() {
        var todoNodes = this.props.data.map(function (todo) {
            return <Todo status={todo.status}>{todo.name}</Todo>;
        });
        return (
                <div className="todoList">
                {todoNodes}
            </div>
        );
    }
});
var TodoForm = React.createClass({
    handleSubmit: function() {
        var status = this.refs.status.getDOMNode().value.trim();
        var name = this.refs.name.getDOMNode().value.trim();
        this.props.onTodoSubmit({status: status, name: name});
        this.refs.status.getDOMNode().value = '';
        this.refs.name.getDOMNode().value = '';
        return false;
    },
    render: function() {
        return (
                <form className="todoForm" onSubmit={this.handleSubmit}>
                <input type="text" placeholder="status" ref="status" />
                <input type="text" placeholder="name" ref="name"/>
                <input type="submit" value="add" />
                </form>
        );
    }
});
var TodoContainer = React.createClass({
    getInitialState: function() {
        return {data: []};
    },
    loadTodosFromServer: function() {
        $.ajax({
            url: this.props.url,
            dataType: 'json',
            success: function(data) {
                this.setState({data: data});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },
    handleTodoSubmit: function(todo) {
        $.ajax({
            url: this.props.url,
            dataType: 'json',
            type: 'POST',
            data: todo,
            success: function(data) {
                this.setState({data: data});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },
    componentWillMount: function() {
        this.loadTodosFromServer();
        setInterval(this.loadTodosFromServer, this.props.pollInterval);
    },
    render: function() {
        return (
                <div className="todoContainer">
                <h1>TODO yer stuffz</h1>
                <TodoList data={this.state.data}/>
                <TodoForm onTodoSubmit={this.handleTodoSubmit} />
            </div>
        );
    }
});
React.renderComponent(
        <TodoContainer url="todo.json" pollInterval={2000} />,
    document.getElementById("content")
);
