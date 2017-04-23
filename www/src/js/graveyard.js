let svg = d3.select('svg');

function updateData(data) {

  let rect = svg.selectAll("rect")
    .data(data);

  rect.enter().append("rect")
    .attr('x', function(d) { return d.coords[0]; })
    .attr('width', function(d) { return d.coords[2] - d.coords[0]; })
    .attr('y', function(d) { return d.coords[1]; })
    .attr('height', function(d) { return d.coords[3] - d.coords[1]; })
    .attr('fill', function(d) {
      if (d.status == 'dead') {
        return '#f00';
      } else {
        return '#333';
      }
    })
    .on('mouseover', function(d) {
      console.log(d)
    });
}

d3.json("data/test.json", function(data) {
  updateData(data);
});
