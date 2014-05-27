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
var UserTime = React.createClass({
    getInitialState: function() {
        return {showtime: true};
    },
    onEditButton: function() {
        this.setState({showtime: false});
    },
    onEdit: function() {
        this.props.onEdit(
            {goalminutes: this.refs.goalminutes.getDOMNode().value.trim()}
        );
        this.setState({showtime: true});
    },
    render: function() {
        var time = function(that) {
            return (
                    <span>Goal: {that.props.goalminutes} <button onClick={that.onEditButton}>edit</button> minutes a day.</span>
            );
        };
        var edit = function(that) {
            return (
                    <span>Goal: <form>
                    <input type="text" ref="goalminutes" defaultValue={that.props.goalminutes} />
                    <button onClick={that.onEdit}>edit</button>
                    </form> minutes a day.</span>
            );
        };
        if (this.state.showtime) {
            return time(this);
        } else {
            return edit(this);
        }
    }
});
var User = React.createClass({
    getInitialState: function() {
        return {person: {name: "nope", goalminutes: 0}};
    },
    handleEditTime: function(time) {
        $.ajax({
            url: "person/time/edit",
            dataType: 'json',
            type: 'POST',
            data: {session: getSession(), goalminutes: time.goalminutes},
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },
    handleUserSubmit: function(person) {
        $.ajax({
            url: "setusername",
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
            url: "person",
            dataType: 'json',
            type: 'POST',
            data: {session: getSession()},
            success: function(person) {
                if (person.name === "") {
                    this.setState({person: {name: "no", goalminutes: 0}});
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
                <h1>
                YOUR NAME IS {this.state.person.name}. <UserTime goalminutes={this.state.person.goalminutes} onEdit={this.handleEditTime} />
                </h1>
                <UserForm onUserSubmit={this.handleUserSubmit} />
            </div>
        );
    }
});
var Todo = React.createClass({
    getInitialState: function() {
        return {showtask: true};
    },
    onCancel: function() {
        this.setState({showtask: true});
        return false;
    },
    onEdit: function() {
        this.props.onEdit({id: this.props.id,
                           name: this.refs.name.getDOMNode().value.trim(),
                           status: this.refs.status.getDOMNode().value.trim(),
                           description: this.refs.description.getDOMNode().value.trim()});
        this.setState({showtask: true});
        return false;
    },
    onEditForm: function() {
        this.setState({
            showtask: false
                      });
        return false;
    },
    onDelete: function() {
        this.props.onDelete({id: this.props.id});
        return false;
    },
    render: function() {
        var form = function(that) {
            return (
                    <form>
                    <input type="text" defaultValue={that.props.status} ref="status"/>
                    <input type="text" defaultValue={that.props.name} ref="name" />
                    <input type="text" defaultValue={that.props.description} ref="description" />
                    <button onClick={that.onEdit}>edit</button>
                    <button onClick={that.onCancel}>cancel</button>
                    </form>
        )};
        var task = function(that) { return (
                <div className="todo">
                <h2 className="status">
                {that.props.status}
            </h2>
                <p>{that.props.name}</p>
                <p>{that.props.description}</p>
                <form onSubmit={that.onEditForm}>
                <button>edit</button>
                </form>
                <form onSubmit={that.onDelete}>
                <button>delete</button>
                </form>
                </div>
        )};
        if (this.state.showtask) {
            return task(this);
        } else {
            return form(this);
        }
    }
});
var TodoList = React.createClass({
    handleEdit: function(task) {
        $.ajax({
            url: "task/edit",
            dataType: 'json',
            type: 'POST',
            data: {session: getSession(), task: task},
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },
    handleDelete: function(id) {
        $.ajax({
            url: "task/delete",
            dataType: 'json',
            type: 'POST',
            data: {session: getSession(), id: id.id},
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },
    render: function() {
        var that = this;
        var todoNodes = this.props.data.map(function (todo) {
            return (
                    <Todo
                onEdit={that.handleEdit}
                onDelete={that.handleDelete}
                status={todo.status}
                name={todo.name}
                id={todo.id}
                description={todo.description}
                    />
            );
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
        var description = this.refs.description.getDOMNode().value.trim();
        this.props.onTodoSubmit({status: status, name: name, description: description});
        this.refs.status.getDOMNode().value = '';
        this.refs.name.getDOMNode().value = '';
        this.refs.description.getDOMNode().value = '';
        return false;
    },
    render: function() {
        return (
                <form className="todoForm" onSubmit={this.handleSubmit}>
                <input type="text" placeholder="status" ref="status" />
                <input type="text" placeholder="name" ref="name"/>
                <input type="text" placeholder="description" ref="description" />
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
            url: "tasks",
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
            url: "addtask",
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
        <User pollInterval={2000} />,
    document.getElementById("user")
);
React.renderComponent(
        <TodoContainer pollInterval={2000} />,
    document.getElementById("tasks")
);
