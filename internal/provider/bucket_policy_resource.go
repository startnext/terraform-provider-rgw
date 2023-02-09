package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = &BucketPolicyResource{}

func NewBucketPolicyResource() resource.Resource {
	return &BucketPolicyResource{}
}

type BucketPolicyResource struct {
	client *RgwClient
}

type BucketPolicyResourceModel struct {
	Id     types.String `tfsdk:"id"`
	Bucket types.String `tfsdk:"bucket"`
	Policy types.String `tfsdk:"policy"`
}

func (r *BucketPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket_policy"
}

func (r *BucketPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Bucket Policy in Ceph RGW",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bucket": schema.StringAttribute{
				MarkdownDescription: "Bucket Name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy": schema.StringAttribute{
				MarkdownDescription: "Bucket Policy",
				Required:            true,
			},
		},
	}
}

func (r *BucketPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*RgwClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *RgwClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *BucketPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var data *BucketPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configure PutBucketPolicy
	s3req := &s3.PutBucketPolicyInput{
		Bucket: aws.String(data.Bucket.ValueString()),
		Policy: aws.String(data.Policy.ValueString()),
	}

	// PutBucketPolicy
	_, err := r.client.S3.PutBucketPolicy(ctx, s3req)
	if err != nil {
		resp.Diagnostics.AddError("could not create bucket policy", err.Error())
		return
	}

	// use bucket name as resource id
	data.Id = types.StringValue(*s3req.Bucket)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BucketPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var data *BucketPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create GetBucketPolicy Request
	s3req := &s3.GetBucketPolicyInput{
		Bucket: aws.String(data.Bucket.ValueString()),
	}

	s3res, err := r.client.S3.GetBucketPolicy(ctx, s3req)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			switch ae.ErrorCode() {
			case "403":
				resp.Diagnostics.AddError("acces denied", "If you are using an identity other than the root user of the Amazon Web Services account that owns the bucket, the calling identity must have the GetBucketPolicy permissions on the specified bucket and belong to the bucket owner's account in order to use this operation")
				return
			case "405":
				resp.Diagnostics.AddError("wrong identity", "If you have the correct permissions, but you're not using an identity that belongs to the bucket owner's account, Amazon S3 returns a 405 Method Not Allowed error.")
				return
			}
		}
		resp.Diagnostics.AddError("could not get bucket policy", err.Error())
		return
	}

	data.Policy = types.StringValue(*s3res.Policy)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BucketPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var data *BucketPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configure PutBucketPolicy
	s3req := &s3.PutBucketPolicyInput{
		Bucket: aws.String(data.Bucket.ValueString()),
		Policy: aws.String(data.Policy.ValueString()),
	}

	// PutBucketPolicy
	_, err := r.client.S3.PutBucketPolicy(ctx, s3req)
	if err != nil {
		resp.Diagnostics.AddError("could not modify bucket policy", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BucketPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Read Terraform prior state data into the model
	var data *BucketPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	s3req := &s3.DeleteBucketPolicyInput{
		Bucket: aws.String(data.Bucket.ValueString()),
	}

	_, err := r.client.S3.DeleteBucketPolicy(ctx, s3req)
	if err != nil {
		resp.Diagnostics.AddError("could not delete bucket policy", err.Error())
		return
	}
}
