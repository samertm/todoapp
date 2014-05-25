/** @jsx React.DOM */
var UserForm = React.createClass({
    handleSubmit: function() {
        var name = this.refs.name.getDOMNode().value.trim();
        this.props.onUserSubmit({name: name});
        this.refs.name.getDOMNode().value = '';
        return false;
    },
    render: function() {
        return (
                <form className="userForm" onSubmit={this.handleSubmit}>
                <input type="text" placeholder="name" ref="name" />
                <input type="submit" value="set" />
                </form>
        );
    }
});
var User = React.createClass({
    getInitialState: function() {
        return {person: {name: "nope"}};
    },
    handleUserSubmit: function(person) {
        $.ajax({
            url: this.props.setusername,
            dataType: 'json',
            type: 'POST',
            data: {session: getSession(), name: person.name},
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },
    loadPersonFromServer: function() {
        $.ajax({
            url: this.props.person,
            dataType: 'json',
            type: 'POST',
            data: {session: getSession()},
            success: function(person) {
                if (person.name === "") {
                    this.setState({person: {name: "no"}});
                } else {
                    this.setState({person: person});
                }
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },    
    componentWillMount: function() {
        this.loadPersonFromServer();
        setInterval(this.loadPersonFromServer, this.props.pollInterval);
    },
    render: function() {
        return (
                <div className="user">
                <h1>YOUR NAME IS {this.state.person.name}</h1>
                <UserForm onUserSubmit={this.handleUserSubmit} />
            </div>
        );
    }
});
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
            url: this.props.tasks,
            dataType: 'json',
            type: 'POST',
            data: {session: getSession()},
            success: function(todos) {
                if (todos !== null) {
                    this.setState({data: todos});
                } else {
                    this.setState({data: []});
                }
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },
    handleTodoSubmit: function(todo) {
        $.ajax({
            url: this.props.addtask,
            dataType: 'json',
            type: 'POST',
            data: {session: getSession(), todo: todo},
            // success: function(data) {
            //     this.setState({data: data});
            // }.bind(this),
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
        <User
    person="person"
    setusername="setusername"
    pollInterval={2000}
        />,
    document.getElementById("user")
);
React.renderComponent(
        <TodoContainer
    addtask="addtask"
    tasks="tasks"
    pollInterval={2000}
        />,
    document.getElementById("tasks")
);
