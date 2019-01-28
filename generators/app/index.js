"use strict";
const Generator = require("yeoman-generator");
const chalk = require("chalk");
const yosay = require("yosay");

module.exports = class extends Generator {
  prompting() {
    // Have Yeoman greet the user.
    this.log(
      yosay(
        `Welcome to the superb ${chalk.red(
          "generator-gin-auth-swag"
        )} generator!`
      )
    );

    const prompts = [
      {
        type: "input",
        name: "myAppPath",
        message:
          "What is the root path for your project (ex: github.com/sunjiaying/generator-gin-api)",
        default: process.cwd().replace(process.env.GOPATH + "/src/", "")
      }
    ];

    return this.prompt(prompts).then(props => {
      // To access props later use this.props.someAnswer;
      this.props = props;
      this.log("app name", props.myAppPath);
      this.myAppPath = props.myAppPath;
    });
  }

  writing() {
    // This.fs.copy(
    //   this.templatePath("dummyfile.txt"),
    //   this.destinationPath("dummyfile.txt")
    // );

    var tmplContext = {
      myAppPath: this.myAppPath
    };

    this.fs.copyTpl(
      this.templatePath("."),
      this.destinationPath("."),
      tmplContext
    );
  }

  install() {
    // This.installDependencies();
  }
};
