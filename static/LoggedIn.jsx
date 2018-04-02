var LoggedIn = React.createClass({
  getInitialState: function() {
    return {
      products: []
    }
  },
  render: function() {
    return (
      <div className="col-lg-12">
        <span className="pull-right"><a onClick={this.logout}>Log out</a></span>
        <h2>Welcome to EconoGopher</h2>
        <p>No better place for you to sit back relax and arm chair economic.</p>
        <div className="row">

        {this.state.products.map(function(product, i){
          return <Product key={i} product={product} />
        })}
        </div>
      </div>);
  }
});
