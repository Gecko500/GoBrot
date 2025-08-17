#version 400
out vec4 frag_colour;

uniform vec2 u_resolution; // window size
uniform double u_zoom;      // zoom factor
uniform dvec2 u_offset;     // pan offset
uniform int max_iter;      // maximum iterations

void main() {
    // Map pixel to complex plane
    dvec2 c = (dvec2(gl_FragCoord.xy) - 0.5 * u_resolution) / u_resolution.y * u_zoom + u_offset;
    dvec2 z = dvec2(0.0);
    // Early bailout for points inside the main cardioid or period-2 bulb
    double x = c.x;
    double y = c.y;
    // Main cardioid check
    double q = (x - 0.25)*(x - 0.25) + y*y;
    if (q*(q + (x - 0.25)) < 0.25*y*y) {
        frag_colour = vec4(0.0, 0.0, 0.0, 1.0);
        return;
    }
    // Period-2 bulb check
    if ((x + 1.0)*(x + 1.0) + y*y < 0.0625) {
        frag_colour = vec4(0.0, 0.0, 0.0, 1.0);
        return;
    }
    
    int i;
    for (i = 0; i < max_iter; i++) {
        if (dot(z, z) > 4.0) break;
        z = dvec2(z.x*z.x - z.y*z.y, 2.0*z.x*z.y) + c;
    }
    if (i == max_iter) {
        frag_colour = vec4(0.0, 0.0, 0.0, 1.0); // inside the set
    } else {
        float r = (i * 2 % 100) / 100.0;
        float g = (i * 6 % 100) / 100.0;
        float b = (i * 5 % 100) / 100.0;
        frag_colour = vec4(r, g, b, 1.0);
    }
}