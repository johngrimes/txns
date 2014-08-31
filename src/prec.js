/** @jsx React.DOM */

var AccountPage = React.createClass({
  handleSubmitSuccess: function() {
    this.refs.txnTable.refreshTxns();
    this.refs.txnUploadForm.render();
  },

  render: function() {
    return (
      <section>
        <h2>Everyday Account</h2>
        <TxnTable ref="txnTable" />
        <h5>Upload a statement</h5>
        <TxnUploadForm ref="txnUploadForm" onSubmitSuccess={this.handleSubmitSuccess} />
      </section>
    );
  }
});

var TxnTable = React.createClass({
  getInitialState: function() {
    return { txns: [] };
  },

  getDefaultProps: function() {
    return {
      txnsURL: 'http://localhost:4000/txns',
      pollFrequency: 10000
    };
  },

  componentDidMount: function() {
    this.refreshTxns();
    setInterval(this.refreshTxns, this.props.pollFrequency);
  },

  refreshTxns: function() {
    $.getJSON(this.props.txnsURL, function(data) {
      this.setState({ txns: data });
    }.bind(this));
  },

  render: function() {
    var txns = this.state.txns;
    return (
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Description</th>
            <th>Debit</th>
            <th>Credit</th>
            <th>Balance</th>
          </tr>
        </thead>
        <tbody>
          {txns.map(function(txn) {
            return <Txn key={txn["id"]} txn={txn} />;
          })}
        </tbody>
      </table>
    );
  }
});

var Txn = React.createClass({
  render: function() {
    var txn = this.props.txn;
    return (
      <tr>
        <td>{txn.date}</td>
        <td>{txn.description}</td>
        <td>{txn.debitCents}</td>
        <td>{txn.creditCents}</td>
        <td>{txn.balanceCents}</td>
      </tr>
    );
  }
});

var TxnUploadForm = React.createClass({
  getDefaultProps: function() {
    return {
      action: 'http://localhost:4000/txns'
    };
  },

  handleSubmit: function() {
    formData = new FormData();
    fileInput = this.refs.fileInput.getDOMNode();
    formData.append(fileInput.name, fileInput.files[0]);
    $.ajax({
      url: this.props.action,
      type: 'POST',
      contentType: false,
      processData: false,
      data: formData,
      success: function(data) {
        this.props.onSubmitSuccess();
      }.bind(this),
      error: function(xhr, err) {
        //TODO: Give error message back to user.
        console.log(err);
      }
    });
    return false;
  },

  render: function() {
    return (
      <form onSubmit={this.handleSubmit} action={this.props.action} method="POST" encType="multipart/form-data">
        <fieldset>
          <input type="file" ref="fileInput" name="txn_file" />
          <input type="submit" />
        </fieldset>
      </form>
    );
  }
});

React.renderComponent(
  <AccountPage />,
  document.getElementById('account')
);
