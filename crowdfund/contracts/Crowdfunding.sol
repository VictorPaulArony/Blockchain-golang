pragma solidity ^0.8.0;

contract Crowdfunding {
    struct Campaign {
        address payable creator;
        uint goal;
        uint deadline;
        uint fundsRaised;
        bool withdrawn;
    }

    mapping(uint => Campaign) public campaigns;
    uint public campaignCount;

    function createCampaign(uint _goal, uint _duration) public {
        campaignCount++;
        campaigns[campaignCount] = Campaign(
            payable(msg.sender),
            _goal,
            block.timestamp + _duration,
            0,
            false
        );
    }

    function contribute(uint _campaignId) public payable {
        Campaign storage campaign = campaigns[_campaignId];
        require(block.timestamp < campaign.deadline, "Campaign expired");
        require(msg.value > 0, "Contribution must be greater than 0");

        campaign.fundsRaised += msg.value;
    }

    function withdrawFunds(uint _campaignId) public {
        Campaign storage campaign = campaigns[_campaignId];
        require(msg.sender == campaign.creator, "Only creator can withdraw");
        require(block.timestamp >= campaign.deadline, "Campaign not ended");
        require(campaign.fundsRaised >= campaign.goal, "Goal not reached");
        require(!campaign.withdrawn, "Funds already withdrawn");

        campaign.withdrawn = true;
        campaign.creator.transfer(campaign.fundsRaised);
    }
}
