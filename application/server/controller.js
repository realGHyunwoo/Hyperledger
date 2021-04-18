var sdk = require('./sdk.js');
module.exports = function(app){
  app.get('/api/getWallet', function (req, res) {
    var walletid = req.query.walletid;
    let args = [walletid];
    sdk.send(false, 'getWallet', args, res);
  });

  app.get('/api/getCertificate', function(req, res){
    var certkey = req.query.certkey;
    let args = [certkey];
    sdk.send(false, 'getCertificate', args, res);
  });

  app.get('/api/setJobPosting', function(req, res){
    var field = req.query.field;
    var name = req.query.name;
    var psword = req.query.psword;
    var conditions = req.query.conditions;
    var pay = req.query.pay;
    var endDate = req.query.endDate;
    let args = [field, name, psword, conditions, pay, endDate];
    sdk.send(true, 'setJobPosting', args, res);
  });

  app.get('/api/setRating', function(req, res){
    var companykey = req.query.companykey;
    var psword = req.query.psword;
    var rating = req.query.rating;
    var jobpostingkey = req.query.jobpostingkey;
    var freelancerkey = req.query.freelancerkey;
    let args = [companykey, psword, rating, jobpostingkey, freelancerkey];
    sdk.send(true, 'setRating', args, res);
  });

  app.get('/api/setApply', function(req, res){
    var dockey = req.query.dockey;
    var psword = req.query.psword;
    var freelancerkey = req.query.freelancerkey;
    let args = [dockey,psword, freelancerkey];
    sdk.send(true, 'setApply', args, res);
  });

  app.get('/api/getJobPosting', function(req, res){
    let args = [];
    sdk.send(false, 'getJobPosting', args, res);
  });
  
}
