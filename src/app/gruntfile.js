module.exports = function(grunt) {
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        concat: {
            options: {
                separator: ';'
            },
            dist: {
                src: ['frontend/**/*.js'],
                dest: 'dist/<%= pkg.name %>.js'
            }
        },
        uglify: {
            options: {
                banner: '/*! <%= pkg.name %> <%= grunt.template.today("yyyy-mm-dd") %> */\n' //添加banner
            },
            buildall: {
                options: {
                    mangle: true,
                    compress: {
                        drop_console: true
                    },
                    report: "min" //输出压缩率，可选的值有 false(不输出信息)，gzip
                },
                files: [{
                    expand: true,
                    cwd: 'frontend', //js目录下
                    src: '**/*.js', //所有js文件
                    dest: 'dist' //输出到此目录下
                }]
            }
        },
        qunit: {
            all: {
                options: {
                    urls: [
                        'http://localhost:8888/index.html',
                    ]
                }
            }
        },
        connect: {
            server: {
                options: {
                    port: 8888,
                    base: '.'
                }
            }
        },
        jshint: {
            files: ['frontend/**/*.js'],
            options: {
                //覆盖默认配置
                globals: {
                    jQuery: true,
                    console: true,
                    module: true,
                    document: true
                }
            }
        },
        watch: {
            scripts: {
                files: ['frontend/**/*.js'],
                tasks: ['default'],
                options: {
                    spawn: true,
                    interrupt: true
                }
            }
        }
    });

    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-jshint');
    grunt.loadNpmTasks('grunt-contrib-qunit');
    grunt.loadNpmTasks('grunt-contrib-watch');
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-connect');

    grunt.registerTask('minall', ['uglify:buildall']);
    grunt.registerTask('unittest', ['connect', 'qunit']);
    grunt.registerTask('test', ['jshint', 'qunit']);
    grunt.registerTask('full', ['jshint', 'qunit', 'concat', 'uglify']);
    grunt.registerTask('default', ['concat', 'uglify']);
};

