var svg = d3.select('svg');
var data = {};
var rect = null;

var filter = '200';

let colorScale = d3.scaleLinear().domain([1,8000])
      .interpolate(d3.interpolateHcl)
      .range([d3.rgb("#f00"), d3.rgb('#0f0')]);

function test() {
  filter = '404';
  updateData(data);
}

function updateData(data) {

  rect = svg.selectAll("rect")
    .data(data, function(d) { return d.title; });

  rect.exit()
    .remove();

  rect.enter()
    .append('rect')
    .merge(rect)
      .attr('x', function(d) { return d.coords.split(",")[0]; })
      .attr('x', function(d) { return d.coords.split(",")[0]; })
      .attr('width', function(d) { return d.coords.split(",")[2] - d.coords.split(",")[0]; })
      .attr('y', function(d) { return d.coords.split(",")[1]; })
      .attr('height', function(d) { return d.coords.split(",")[3] - d.coords.split(",")[1]; })
      .attr('fill', function(d) {
        width = d.coords.split(",")[2] - d.coords.split(",")[0];
        height = d.coords.split(",")[3] - d.coords.split(",")[1];
        size = width * height;
        return colorScale(size);
      })
      .on('mouseover', function(d) {
        console.log(d)
      });
}

d3.json("data/data.json", function(d) {
  console.log(d);
  data = d;
  updateData(d);
});
